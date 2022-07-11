package core

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/feimingxliu/quicksearch/internal/config"
	_ "github.com/feimingxliu/quicksearch/internal/pkg/analyzer"
	"github.com/feimingxliu/quicksearch/internal/pkg/analyzer/jieba"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/uuid"
	"os"
	"path"
	"sync"
	"time"
)

type Index struct {
	UID            string        `json:"uid"`
	Name           string        `json:"name"`
	Mapping        *IndexMapping `json:"mapping"`
	DocNum         uint64        `json:"doc_num"`          // number of docs
	StorageSize    uint64        `json:"storage_size"`     // bytes on disk
	NumberOfShards int           `json:"number_of_shards"` // number of shards
	Shards         []*IndexShard `json:"shards"`
	CreateAt       time.Time     `json:"create_at"`
	UpdateAt       time.Time     `json:"update_at"`
	mu             sync.RWMutex
}

type options struct {
	name        string
	mapping     *IndexMapping
	numOfShards int
}

type Option func(*options)

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithIndexMapping(m *IndexMapping) Option {
	return func(o *options) {
		o.mapping = m
	}
}

func WithShards(num int) Option {
	return func(o *options) {
		o.numOfShards = num
	}
}

// NewIndex return an Index, which is opened and appended to engine.indices.
func NewIndex(opts ...Option) (*Index, error) {
	// the default will be replaced by opts
	cfg := &options{}
	for _, opt := range opts {
		opt(cfg)
	}
	if index, err := GetIndex(cfg.name); err == nil && index != nil {
		return index, nil
	}
	uid := uuid.GetXID()
	index := &Index{
		UID:            uid,
		Name:           cfg.name,
		Mapping:        cfg.mapping,
		NumberOfShards: cfg.numOfShards,
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
	}
	if index.NumberOfShards <= 0 {
		index.NumberOfShards = config.Global.Engine.DefaultNumberOfShards
	}
	// open and put into engine.indices
	if err := index.Open(); err != nil {
		return nil, err
	}
	return index, nil
}

func (index *Index) SetMapping(mapping interface{}) error {
	if mapping == nil {
		return nil
	}
	switch m := mapping.(type) {
	case *IndexMapping:
		index.Mapping = m
	case IndexMapping:
		index.Mapping = &m
	case map[string]interface{}:
		if mp, err := BuildIndexMappingFromMap(m); err != nil {
			return err
		} else {
			index.Mapping = mp
		}
	case *map[string]interface{}:
		if mp, err := BuildIndexMappingFromMap(*m); err != nil {
			return err
		} else {
			index.Mapping = mp
		}
	default:
		return errors.ErrInvalidMapping
	}
	return nil
}

// ListIndices list all managed indices, note: not opened!
func ListIndices() ([]*Index, error) {
	data, err := engine.meta.List()
	if err != nil {
		return nil, err
	}
	indexes := make([]*Index, 0, len(data))
	for _, d := range data {
		idx := new(Index)
		err = json.Unmarshal(d, idx)
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, idx)
	}
	return indexes, nil
}

// GetIndex firstly search in mem, than find in db, in err == nil and index != nil, it's ready to use(opened).
func GetIndex(name string) (*Index, error) {
	if index := engine.getIndex(name); index != nil {
		// the index got from engine cache already opens
		return index, nil
	}
	b, err := engine.meta.Get(name)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			return nil, errors.ErrIndexNotFound
		}
		return nil, err
	}
	index := new(Index)
	err = json.Unmarshal(b, index)
	if err != nil {
		return nil, err
	}
	if err = index.Open(); err != nil {
		return nil, err
	}
	return index, nil
}

// Open open the index, append it to the engine.indices.
func (index *Index) Open() error {
	index.mu.Lock()
	if index.Shards == nil {
		index.Shards = make([]*IndexShard, 0, index.NumberOfShards)
		for i := 0; i < index.NumberOfShards; i++ {
			mapping, err := buildIndexMapping(index.Mapping)
			if err != nil {
				return err

			}
			indexer, err := bleve.New(index.shardDir(i), mapping)
			if err != nil {
				return err
			}
			shard := &IndexShard{
				ID:      i,
				Indexer: indexer,
			}
			index.Shards = append(index.Shards, shard)
		}
	} else {
		for _, shard := range index.Shards {
			if shard.Indexer != nil {
				continue
			}
			indexer, err := bleve.Open(index.shardDir(shard.ID))
			if err != nil {
				index.mu.Unlock()
				return err
			}
			shard.Indexer = indexer
		}
	}
	engine.addIndex(index)
	index.mu.Unlock()
	// update metadata after open
	if err := index.UpdateMetadata(); err != nil {
		return err
	}
	return nil
}

// Close closes index and release the related resource, including remove from engine.indices.
func (index *Index) Close() error {
	// update metadata before close
	if err := index.UpdateMetadata(); err != nil {
		return err
	}
	index.mu.Lock()
	engine.removeIndex(index)
	for _, shard := range index.Shards {
		if shard.Indexer == nil {
			continue
		}
		// cleanup cgo allocated heap memory
		if az := shard.Indexer.Mapping().AnalyzerNamed("gojieba"); az != nil {
			az.Tokenizer.(*jieba.JiebaTokenizer).Free()
		}
		if err := shard.Indexer.Close(); err != nil {
			index.mu.Unlock()
			return err
		}
		shard.Indexer = nil
	}
	index.mu.Unlock()
	return nil
}

// Delete close the index, remove from engine.indices, delete all index files, delete index metadata.
func (index *Index) Delete() error {
	// close
	if err := index.Close(); err != nil {
		return err
	}
	index.mu.Lock()
	defer index.mu.Unlock()
	// delete metadata.
	if err := engine.meta.Delete(index.Name); err != nil {
		return err
	}
	// remove all files.
	if err := os.RemoveAll(index.dir()); err != nil {
		return err
	}
	return nil
}

// Clone clones the entire index to a new index.
func (index *Index) Clone(name string) error {
	// check if cloned index is valid
	if _, err := GetIndex(name); err == nil {
		return errors.ErrIndexAlreadyExists
	} else {
		if err != errors.ErrIndexNotFound {
			return err
		}
	}
	uid := uuid.GetXID()
	clone := &Index{
		UID:            uid,
		Name:           name,
		Mapping:        index.Mapping,
		DocNum:         index.DocNum,
		StorageSize:    index.StorageSize,
		NumberOfShards: index.NumberOfShards,
		Shards:         make([]*IndexShard, 0, index.NumberOfShards),
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		mu:             sync.RWMutex{},
	}
	// clone all shards.
	// open the index first if it is closed.
	if err := index.Open(); err != nil {
		return err
	}
	index.mu.RLock()
	for _, shard := range index.Shards {
		copyableIndex, ok := shard.Indexer.(bleve.IndexCopyable)
		if !ok {
			index.mu.RUnlock()
			return errors.ErrIndexCloneNotSupported
		}
		if err := copyableIndex.CopyTo(bleve.FileSystemDirectory(clone.shardDir(shard.ID))); err != nil {
			return err
		}
		clone.Shards = append(clone.Shards, &IndexShard{
			ID:          shard.ID,
			DocNum:      shard.DocNum,
			StorageSize: shard.StorageSize,
			Indexer:     nil,
		})
	}
	index.mu.RUnlock()
	// open the cloned index.
	if err := clone.Open(); err != nil {
		return err
	}
	// write metadata.
	if err := clone.UpdateMetadata(); err != nil {
		return err
	}
	return nil
}

// UpdateMetadata writes the index metadata to db.
func (index *Index) UpdateMetadata() error {
	var totalDocNum, totalSize uint64
	// update docNum and storageSize
	for i := 0; i < index.NumberOfShards; i++ {
		index.UpdateMetadataByShard(i)
	}
	index.mu.RLock()
	for i := 0; i < index.NumberOfShards; i++ {
		totalDocNum += index.Shards[i].DocNum
		totalSize += index.Shards[i].StorageSize
	}
	if totalDocNum > 0 && totalSize > 0 {
		index.DocNum = totalDocNum
		index.StorageSize = totalSize
	}
	b, _ := json.Marshal(index)
	index.mu.RUnlock()
	return engine.meta.Set(index.Name, b)
}

func (index *Index) UpdateMetadataByShard(n int) {
	index.mu.RLock()
	shard := index.Shards[n]
	index.mu.RUnlock()
	if shard.Indexer == nil {
		return
	}
	docNum, _ := shard.Indexer.DocCount()
	var storageSize uint64
	if n, ok := shard.Indexer.StatsMap()["CurOnDiskBytes"].(uint64); ok {
		storageSize = n
	}
	if docNum > 0 {
		shard.DocNum = docNum
	}
	if storageSize > 0 {
		shard.StorageSize = storageSize
	}
}

// returns the index storage dir
func (index *Index) dir() string {
	return path.Join(config.Global.Storage.DataDir, "indices", index.Name)
}

// returns the index shard storage dir
func (index *Index) shardDir(shard int) string {
	return path.Join(index.dir(), fmt.Sprintf("%s_%d", index.UID, shard))
}

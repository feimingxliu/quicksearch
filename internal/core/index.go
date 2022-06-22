package core

import (
	pconfig "github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/pkg/storager"
	ptokenizer "github.com/feimingxliu/quicksearch/internal/pkg/tokenizer"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/base64"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Index struct {
	Name          string    `json:"name"`
	StorageType   string    `json:"storage_type"`
	TokenizerType string    `json:"tokenizer_type"`
	DocNum        uint64    `json:"doc_num"`
	DocTimeMin    int64     `json:"doc_time_min"` //indexed doc's min @timestamp (ns)
	DocTimeMax    int64     `json:"doc_time_max"` //indexed doc's max @timestamp (ns)
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at"`
	StorePath     string    `json:"store_path"`    //document storage file path
	InvertedPath  string    `json:"inverted_path"` //inverted index storage file path
	rwMutex       sync.RWMutex
	open          bool
	store         storager.Storager //for document storage
	inverted      storager.Storager //for inverted index storage
	tokenizer     ptokenizer.Tokenizer
}

type options struct {
	name          string
	storageType   string
	tokenizerType string
}

type Option func(*options)

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithStorage(s string) Option {
	return func(o *options) {
		o.storageType = s
	}
}

func WithTokenizer(t string) Option {
	return func(o *options) {
		o.tokenizerType = t
	}
}

//NewIndex return an Index, which is opened and appended to Indices.
func NewIndex(opts ...Option) *Index {
	config := new(options)
	for _, opt := range opts {
		opt(config)
	}
	if index, err := GetIndex(config.name); err == nil && index != nil {
		return index
	}
	index := &Index{
		Name:          config.name,
		StorageType:   strings.ToLower(config.storageType),
		TokenizerType: strings.ToLower(config.tokenizerType),
		StorePath:     path.Join(pconfig.Global.Storage.DataDir, "indices", base64.Encode([]byte(config.name))),
		InvertedPath:  path.Join(pconfig.Global.Storage.DataDir, "inverted", base64.Encode([]byte(config.name))),
		CreateAt:      time.Now(),
		UpdateAt:      time.Now(),
	}
	switch index.StorageType {
	case "bolt":
		index.store = storager.NewStorager(
			storager.Bolt, index.StorePath,
		)
		index.inverted = storager.NewStorager(
			storager.Bolt, index.InvertedPath,
		)
	default:
		index.store = storager.NewStorager(
			storager.Default, index.StorePath,
		)
		index.inverted = storager.NewStorager(
			storager.Default, index.InvertedPath,
		)
	}
	switch index.TokenizerType {
	case "jieba":
		index.tokenizer = ptokenizer.NewTokenizer(ptokenizer.Jieba)
	default:
		index.tokenizer = ptokenizer.NewTokenizer(ptokenizer.Default)
	}
	index.open = true
	indicesRwMu.Lock()
	Indices[config.name] = index
	indicesRwMu.Unlock()
	//store metadata.
	_ = index.UpdateMetadata()
	return index
}

//ListIndices list all managed indices, note: not opened!
func ListIndices() ([]*Index, error) {
	data, err := meta.List()
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

//GetIndex firstly search in mem, than find in db, in err == nil and index != nil, it's ready to use(opened).
func GetIndex(name string) (*Index, error) {
	indicesRwMu.RLock()
	if index, ok := Indices[name]; ok {
		indicesRwMu.RUnlock()
		return index, nil
	}
	indicesRwMu.RUnlock()
	b, err := meta.Get(name)
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

//Open open the index as well as append to the Indices.
func (index *Index) Open() error {
	{
		index.rwMutex.RLock()
		if index.open {
			index.rwMutex.RUnlock()
			return nil
		}
		index.rwMutex.RUnlock()
	}

	{
		index.rwMutex.Lock()
		switch index.StorageType {
		case "bolt":
			index.store = storager.NewStorager(
				storager.Bolt, index.StorePath,
			)
			index.inverted = storager.NewStorager(
				storager.Bolt, index.InvertedPath,
			)
		default:
			index.store = storager.NewStorager(
				storager.Default, index.StorePath,
			)
			index.inverted = storager.NewStorager(
				storager.Default, index.InvertedPath,
			)
		}
		switch index.TokenizerType {
		case "jieba":
			index.tokenizer = ptokenizer.NewTokenizer(ptokenizer.Jieba)
		default:
			index.tokenizer = ptokenizer.NewTokenizer(ptokenizer.Default)
		}
		index.open = true
		index.rwMutex.Unlock()
	}

	{
		indicesRwMu.Lock()
		Indices[index.Name] = index
		indicesRwMu.Unlock()
	}

	return nil
}

//Close closes index and release the related resource, including remove from Indices.
func (index *Index) Close() error {
	{
		index.rwMutex.RLock()
		if !index.open {
			index.rwMutex.RUnlock()
			return nil
		}
		index.rwMutex.RUnlock()
	}

	{
		index.rwMutex.Lock()
		if err := index.store.Close(); err != nil {
			return err
		}
		if err := index.inverted.Close(); err != nil {
			return err
		}
		index.tokenizer.Close()
		index.open = false
		index.rwMutex.Unlock()
	}

	{
		indicesRwMu.Lock()
		delete(Indices, index.Name)
		indicesRwMu.Unlock()
	}

	return nil
}

//Delete close the index, remove from Indices, delete all docs within it, delete index metadata, remove db file and inverted file.
func (index *Index) Delete() error {
	indicesRwMu.Lock()
	delete(Indices, index.Name)
	indicesRwMu.Unlock()
	index.rwMutex.Lock()
	defer index.rwMutex.Unlock()
	index.open = false
	//close tokenizer.
	index.tokenizer.Close()
	//delete docs.
	if err := index.store.DeleteAll(); err != nil {
		return err
	}
	//delete inverted index.
	if err := index.inverted.DeleteAll(); err != nil {
		return err
	}
	//delete metadata.
	if err := meta.Delete(index.Name); err != nil {
		return err
	}
	//close db.
	if err := index.store.Close(); err != nil {
		return err
	}
	if err := index.inverted.Close(); err != nil {
		return err
	}
	//remove db file.
	if err := os.Remove(index.StorePath); err != nil {
		return err
	}
	if err := os.Remove(index.InvertedPath); err != nil {
		return err
	}
	return nil
}

//Clone clones the entire index to a new index, docs included.
func (index *Index) Clone(name string) error {
	if name == index.Name {
		return errors.ErrCloneIndexSameName
	}
	clone := &Index{
		Name:          name,
		StorageType:   index.StorageType,
		TokenizerType: index.TokenizerType,
		DocNum:        index.DocNum,
		DocTimeMin:    index.DocTimeMin,
		DocTimeMax:    index.DocTimeMax,
		CreateAt:      time.Now(),
		UpdateAt:      time.Now(),
		StorePath:     path.Join(path.Dir(index.StorePath), base64.Encode([]byte(name))),
		InvertedPath:  path.Join(path.Dir(index.InvertedPath), base64.Encode([]byte(name))),
	}
	//clone storage file.
	if err := index.store.CloneDatabase(clone.StorePath); err != nil {
		return err
	}
	if err := index.inverted.CloneDatabase(clone.InvertedPath); err != nil {
		return err
	}
	//write metadata.
	if err := clone.UpdateMetadata(); err != nil {
		return err
	}
	//open the index.
	if err := clone.Open(); err != nil {
		return err
	}
	return nil
}

//SetTimestamp updates DocTimeMin and DocTimeMax.
func (index *Index) SetTimestamp(t int64) {
	if index.DocTimeMin == 0 {
		atomic.StoreInt64(&index.DocTimeMin, t)
	}
	if index.DocTimeMax == 0 {
		atomic.StoreInt64(&index.DocTimeMax, t)
	}
	if t < index.DocTimeMin {
		atomic.StoreInt64(&index.DocTimeMin, t)
	}
	if t > index.DocTimeMax {
		atomic.StoreInt64(&index.DocTimeMax, t)
	}
}

//UpdateMetadata writes the index metadata to db.
func (index *Index) UpdateMetadata() error {
	index.rwMutex.RLock()
	b, err := json.Marshal(index)
	if err != nil {
		return err
	}
	index.rwMutex.RUnlock()
	return meta.Set(index.Name, b)
}

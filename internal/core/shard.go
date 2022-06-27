package core

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/internal/pkg/storager"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"os"
	"path"
)

const (
	DefaultShards = 10 // default number of shards.
)

func NewShards(config *ShardConfig) *Shards {
	if config.NumOfShards <= 0 {
		config.NumOfShards = DefaultShards
	}
	shards := &Shards{
		path:        config.Path,
		indexName:   config.IndexName,
		numOfShards: config.NumOfShards,
		storages:    newStorages(config),
	}
	return shards
}

func newStorages(config *ShardConfig) []storager.Storager {
	storages := make([]storager.Storager, 0, config.NumOfShards)
	for i := 0; i < config.NumOfShards; i++ {
		s := storager.NewStorager(config.StorageType, path.Join(config.Path, fmt.Sprintf("%s_%d", config.IndexName, i)))
		storages = append(storages, s)
	}
	return storages
}

type ShardConfig struct {
	Path        string
	IndexName   string
	StorageType storager.StorageType
	NumOfShards int
}

type Shards struct {
	path        string
	indexName   string
	numOfShards int
	storages    []storager.Storager
}

func (s *Shards) Get(key string) ([]byte, error) {
	//get the shard.
	ss := s.storages[s.getShard(key)]
	return ss.Get(key)
}

func (s *Shards) Set(key string, value []byte) error {
	//get the shard.
	ss := s.storages[s.getShard(key)]
	return ss.Set(key, value)
}

func (s *Shards) Batch(keys []string, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.ErrKeyValueNotMatch
	}
	shardKey := make([][]int, s.numOfShards)
	for idx, key := range keys {
		shard := s.getShard(key)
		shardKey[shard] = append(shardKey[shard], idx)
	}
	for shard, idxes := range shardKey {
		if len(idxes) == 0 {
			continue
		}
		ks, vs := make([]string, 0, len(idxes)), make([][]byte, 0, len(idxes))
		for _, idx := range idxes {
			ks = append(ks, keys[idx])
			vs = append(vs, values[idx])
		}
		if err := s.storages[shard].Batch(ks, vs); err != nil {
			return err
		}
	}
	return nil
}

func (s *Shards) Delete(key string) error {
	//get the shard.
	ss := s.storages[s.getShard(key)]
	return ss.Delete(key)
}

//DeleteAll deletes the index.
func (s *Shards) DeleteAll() error {
	for i := range s.storages {
		if err := s.storages[i].DeleteAll(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Shards) CloneIndex(indexDir string) error {
	if err := os.MkdirAll(indexDir, 0755); err != nil {
		return err
	}
	for i := range s.storages {
		if err := s.storages[i].CloneDatabase(path.Join(indexDir, fmt.Sprintf("%s_%d", path.Base(indexDir), i))); err != nil {
			return err
		}
	}
	return nil
}

func (s *Shards) Close() error {
	for i := range s.storages {
		if err := s.storages[i].Close(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Shards) getShard(key string) int64 {
	return util.BytesModInt([]byte(key), s.numOfShards)
}

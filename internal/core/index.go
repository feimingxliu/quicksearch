package core

import (
	"github.com/feimingxliu/quicksearch/internal/pkg/storager"
	"github.com/feimingxliu/quicksearch/internal/pkg/tokenizer"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"sync"
	"sync/atomic"
	"time"
)

type Index struct {
	Name        string              `json:"name"`
	StorageType string              `json:"storage_type"`
	DocNum      uint64              `json:"doc_num"`
	DocTimeMin  int64               `json:"doc_time_min"`
	DocTimeMax  int64               `json:"doc_time_max"`
	CreateAt    time.Time           `json:"create_at"`
	UpdateAt    time.Time           `json:"update_at"`
	tokenizer   tokenizer.Tokenizer `json:"-"`
	storager    storager.Storager   `json:"-"`
	lock        sync.RWMutex        `json:"-"`
}

func ListIndexes() ([]*Index, error) {
	data, err := db.List("/index/")
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

func GetIndex(name string) (*Index, error) {
	b, err := db.Get(key(name))
	if err != nil {
		return nil, err
	}
	index := new(Index)
	err = json.Unmarshal(b, index)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func DeleteIndex(name string) error {
	return db.DeleteAll(key(name))
}

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

func (index *Index) UpdateMetadata() error {
	b, err := json.Marshal(index)
	if err != nil {
		return err
	}
	return db.Set(key(index.Name), b)
}

func key(name string) string {
	return "/index/" + name
}

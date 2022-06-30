package core

import (
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	DefaultExpiration = 30 * time.Second // time for key expiration in cache
	CleanupInterval   = 1 * time.Minute  // time for clean expired key in cache
	invertedRwMu      sync.RWMutex
	invertedCache     = make(map[*Index]*cache.Cache)
)

//MapKeywordsDoc maps keywords to the doc.
func (index *Index) MapKeywordsDoc(keywords []string, docID string) error {
	if err := index.Open(); err != nil {
		return err
	}
	var keys []string
	var values [][]byte
	c := invertedCache[index]
	for _, keyword := range keywords {
		if len(keyword) == 0 {
			continue
		}
		var ids []string
		// search in cache first
		v, found := c.Get(keyword)
		if found {
			ids = v.([]string)
		} else {
			bids, _ := index.inverted.Get(keyword)
			if bids != nil {
				_ = json.Unmarshal(bids, &ids)
			}
		}
		if slices.ContainsStr(ids, docID) {
			continue
		}
		ids = append(ids, docID)
		// load into cache
		c.Set(keyword, ids, DefaultExpiration)
		bids, err := json.Marshal(ids)
		if err != nil {
			return err
		}
		keys = append(keys, keyword)
		values = append(values, bids)
	}
	err := index.inverted.Batch(keys, values)
	if err != nil {
		return err
	}
	return nil
}

//UnMapKeywordsDoc unmaps keywords to the doc.
func (index *Index) UnMapKeywordsDoc(keywords []string, docID string) error {
	if err := index.Open(); err != nil {
		return err
	}
	var keys []string
	var values [][]byte
	c := invertedCache[index]
	for _, keyword := range keywords {
		if len(keyword) == 0 {
			continue
		}
		var ids []string
		// search in cache first
		v, found := c.Get(keyword)
		if found {
			ids = v.([]string)
		} else {
			bids, _ := index.inverted.Get(keyword)
			if bids != nil {
				_ = json.Unmarshal(bids, &ids)
			}
		}
		ids = slices.RemoveSpecifiedStr(ids, docID)
		// load into cache
		c.Set(keyword, ids, DefaultExpiration)
		bids, err := json.Marshal(ids)
		if err != nil {
			return err
		}
		keys = append(keys, keyword)
		values = append(values, bids)
	}
	err := index.inverted.Batch(keys, values)
	if err != nil {
		return err
	}
	return nil
}

//GetIDsByKeyword retrieves the associated IDs about the keyword.
func (index *Index) GetIDsByKeyword(keyword string) ([]string, error) {
	if err := index.Open(); err != nil {
		return nil, err
	}
	var ids []string
	c := invertedCache[index]
	// search in cache first
	v, found := c.Get(keyword)
	if found {
		ids = v.([]string)
	} else {
		bids, _ := index.inverted.Get(keyword)
		if bids != nil {
			_ = json.Unmarshal(bids, &ids)
		}
	}
	return ids, nil
}

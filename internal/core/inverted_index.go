package core

import (
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
)

type keywordsDoc struct {
	add      bool
	keywords []string
	docID    string
}

func (index *Index) runInvertedIndexWorker() func() {
	toStop := make(chan struct{})
	stopped := make(chan struct{})
	for i := 0; i < index.NumberOfShards; i++ {
		go func(shard int) {
			workCh := index.invertIndexCh[shard]
			cache := index.invertedCache
		loop:
			for {
				select {
				case kd := <-workCh:
					if kd.add {
						values := make([][]byte, 0, len(kd.keywords))
						for _, keyword := range kd.keywords {
							if len(keyword) == 0 {
								continue
							}
							var ids []string
							// search in cache first
							v, found := cache.Get(keyword)
							if found {
								ids = v.([]string)
							} else {
								bids, _ := index.inverted.storages[shard].Get(keyword)
								if bids != nil {
									_ = json.Unmarshal(bids, &ids)
								}
							}
							ids = append(ids, kd.docID)
							// load into cache
							cache.Set(keyword, ids, DefaultExpiration)
							bids, _ := json.Marshal(ids)
							values = append(values, bids)
							_ = index.inverted.storages[shard].Batch(kd.keywords, values)
						}
					} else {
						values := make([][]byte, 0, len(kd.keywords))
						for _, keyword := range kd.keywords {
							if len(keyword) == 0 {
								continue
							}
							var ids []string
							// search in cache first
							v, found := cache.Get(keyword)
							if found {
								ids = v.([]string)
							} else {
								bids, _ := index.inverted.storages[shard].Get(keyword)
								if bids != nil {
									_ = json.Unmarshal(bids, &ids)
								}
							}
							ids = slices.RemoveSpecifiedStr(ids, kd.docID)
							// update cache
							cache.Set(keyword, ids, DefaultExpiration)
							bids, _ := json.Marshal(ids)
							values = append(values, bids)
							_ = index.inverted.storages[shard].Batch(kd.keywords, values)
						}
					}
				case <-toStop:
					stopped <- struct{}{}
					break loop
				}
			}
		}(i)
	}
	return func() {
		close(toStop)
		for i := 0; i < index.NumberOfShards; i++ {
			<-stopped
		}
		close(stopped)
	}
}

// TODO: fix the concurrency bug.
//MapKeywordsDoc maps keywords to the doc.
func (index *Index) MapKeywordsDoc(keywords []string, docID string) error {
	if err := index.Open(); err != nil {
		return err
	}
	shardKeywords := make([][]string, index.NumberOfShards)
	for _, keyword := range keywords {
		shard := index.inverted.getShard(keyword)
		shardKeywords[shard] = append(shardKeywords[shard], keyword)
	}
	for i := range shardKeywords {
		index.invertIndexCh[i] <- &keywordsDoc{
			add:      true,
			keywords: shardKeywords[i],
			docID:    docID,
		}
	}
	return nil
}

//UnMapKeywordsDoc unmaps keywords to the doc.
func (index *Index) UnMapKeywordsDoc(keywords []string, docID string) error {
	if err := index.Open(); err != nil {
		return err
	}
	shardKeywords := make([][]string, index.NumberOfShards)
	for _, keyword := range keywords {
		shard := index.inverted.getShard(keyword)
		shardKeywords[shard] = append(shardKeywords[shard], keyword)
	}
	for i := range shardKeywords {
		index.invertIndexCh[i] <- &keywordsDoc{
			add:      false,
			keywords: shardKeywords[i],
			docID:    docID,
		}
	}
	return nil
}

//GetIDsByKeyword retrieves the associated IDs about the keyword.
func (index *Index) GetIDsByKeyword(keyword string) ([]string, error) {
	if err := index.Open(); err != nil {
		return nil, err
	}
	var ids []string
	c := index.invertedCache
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

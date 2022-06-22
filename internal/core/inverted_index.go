package core

import (
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
)

//TODO: optimize MapKeywordsDoc, it's too slow.

//MapKeywordsDoc maps keywords to the doc.
func (index *Index) MapKeywordsDoc(keywords []string, docID string) error {
	if err := index.Open(); err != nil {
		return err
	}
	keys := make([]string, 0)
	values := make([][]byte, 0)
	for _, keyword := range keywords {
		if len(keyword) == 0 {
			continue
		}
		bids, err := index.inverted.Get(keyword)
		var ids []string
		if err != nil {
			if err == errors.ErrKeyNotFound {
				ids = make([]string, 0)
			} else {
				return err
			}
		}
		if bids != nil {
			err = json.Unmarshal(bids, &ids)
			if err != nil {
				return err
			}
		}
		if slices.ContainsStr(ids, docID) {
			continue
		}
		ids = append(ids, docID)
		bids, err = json.Marshal(ids)
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
	keys := make([]string, 0)
	values := make([][]byte, 0)
	for _, keyword := range keywords {
		if len(keyword) == 0 {
			continue
		}
		bids, err := index.inverted.Get(keyword)
		var ids []string
		if err != nil {
			if err == errors.ErrKeyNotFound {
				//this can never happen.
				continue
			} else {
				return err
			}
		}
		if bids != nil {
			err = json.Unmarshal(bids, &ids)
			if err != nil {
				return err
			}
		}
		ids = slices.RemoveSpecifiedStr(ids, docID)
		bids, err = json.Marshal(ids)
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
	bids, err := index.inverted.Get(keyword)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	err = json.Unmarshal(bids, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

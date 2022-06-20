package core

import (
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
)

var InvertedIndex inverted

type inverted struct{}

func (iv inverted) MapKeywordsDoc(keywords map[string]struct{}, docID string) error {
	for keyword := range keywords {
		bids, err := db.Get(invertedKey(keyword))
		var ids []string
		if err != nil {
			if err == errors.ErrKeyNotFound {
				ids = make([]string, 1)
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
		err = db.Set(invertedKey(keyword), bids)
		if err != nil {
			return err
		}
	}
	return nil
}

func (iv inverted) GetIDsByKeyword(keyword string) ([]string, error) {
	bids, err := db.Get(invertedKey(keyword))
	var ids []string
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bids, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func invertedKey(keyword string) string {
	return "/inverted/" + keyword
}

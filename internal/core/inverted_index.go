package core

import (
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
)

var InvertedIndex inverted

type inverted struct{}

type Pair struct {
	Key string   `json:"keyword"`
	Ids []string `json:"ids"`
}

func (iv inverted) MapKeywordsDoc(keywords []string, docID string) error {
	for _, keyword := range keywords {
		bpair, err := db.Get(invertedKey(keyword))
		var pair Pair
		if err != nil {
			if err == errors.ErrKeyNotFound {
				pair = Pair{
					Key: keyword,
					Ids: make([]string, 0, 1),
				}
			} else {
				return err
			}
		}
		if bpair != nil {
			err = json.Unmarshal(bpair, &pair)
			if err != nil {
				return err
			}
		}
		if slices.ContainsStr(pair.Ids, docID) {
			continue
		}
		pair.Ids = append(pair.Ids, docID)
		bpair, err = json.Marshal(pair)
		if err != nil {
			return err
		}
		err = db.Set(invertedKey(keyword), bpair)
		if err != nil {
			return err
		}
	}
	return nil
}

func (iv inverted) GetIDsByKeyword(keyword string) ([]string, error) {
	bpair, err := db.Get(invertedKey(keyword))
	if err != nil {
		if err == errors.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	var pair Pair
	err = json.Unmarshal(bpair, &pair)
	if err != nil {
		return nil, err
	}
	return pair.Ids, nil
}

func invertedKey(keyword string) string {
	return "/inverted/" + keyword
}

func (iv inverted) listAllKeywordIDs() ([]*Pair, error) {
	b, err := db.List("/inverted/")
	if err != nil {
		return nil, err
	}
	pairs := make([]*Pair, 0, len(b))
	for i := range b {
		pair := new(Pair)
		if err = json.Unmarshal(b[i], pair); err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}
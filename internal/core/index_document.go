package core

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/maps"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
	"strings"
	"time"
)

func (index *Index) IndexDocument(doc *Document) error {
	if doc == nil {
		return nil
	}
	if err := index.Open(); err != nil {
		return err
	}
	//if doc already exists, remove old keyword doc map.
	var update bool
	if bdoc, err := index.store.Get(doc.ID); err == nil {
		//shadowed doc.
		doc := new(Document)
		if err = json.Unmarshal(bdoc, doc); err != nil {
			return err
		}
		if err = index.UnMapKeywordsDoc(doc.KeyWords, doc.ID); err != nil {
			return err
		}
		update = true
	}
	doc.IndexName = index.Name
	doc.Index = index
	//update index metadata.
	index.SetTimestamp(doc.Timestamp.UnixNano())
	index.rwMutex.Lock()
	if !update {
		index.DocNum++
	}
	index.UpdateAt = time.Now()
	index.rwMutex.Unlock()
	//write to db.
	if err := index.UpdateMetadata(); err != nil {
		return err
	}
	//extract tokens and add or update inverted index.
	flatDoc := maps.Flatten(doc.Source)
	keywords := make(map[string]struct{})
	for _, value := range flatDoc {
		//this casts both string and number to string.
		s := fmt.Sprint(value)
		//filter the key char '/' and blank token.
		kws := slices.RemoveEmptyStr(slices.FilterStr(index.tokenizer.Tokenize(s), func(token string) string {
			return strings.TrimSpace(token)
		}))
		for _, token := range kws {
			if _, ok := keywords[token]; !ok {
				keywords[token] = struct{}{}
			}
		}
	}
	doc.KeyWords = make([]string, 0)
	for kw := range keywords {
		doc.KeyWords = append(doc.KeyWords, kw)
	}
	err := index.MapKeywordsDoc(doc.KeyWords, doc.ID)
	if err != nil {
		return err
	}
	//write doc into db.
	b, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = index.store.Set(doc.ID, b)
	if err != nil {
		return err
	}
	return nil
}

func (index *Index) BulkDocuments(docs []*Document) error {
	if len(docs) == 0 {
		return nil
	}
	if err := index.Open(); err != nil {
		return err
	}
	keys := make([]string, len(docs), len(docs))
	values := make([][]byte, len(docs), len(docs))
	totalBulk := len(docs)
	for i, doc := range docs {
		//if doc already exists, remove old keyword doc map.
		if bdoc, err := index.store.Get(doc.ID); err == nil {
			//shadowed doc.
			doc := new(Document)
			if err = json.Unmarshal(bdoc, doc); err != nil {
				return err
			}
			if err = index.UnMapKeywordsDoc(doc.KeyWords, doc.ID); err != nil {
				return err
			}
			totalBulk--
		}
		doc.IndexName = index.Name
		doc.Index = index
		//update index metadata.
		index.SetTimestamp(doc.Timestamp.UnixNano())
		//extract tokens and add or update inverted index.
		flatDoc := maps.Flatten(doc.Source)
		keywords := make(map[string]struct{})
		for _, value := range flatDoc {
			//this casts both string and number to string.
			s := fmt.Sprint(value)
			//filter the key char '/' and blank token.
			kws := slices.RemoveEmptyStr(slices.FilterStr(index.tokenizer.Tokenize(s), func(token string) string {
				return strings.TrimSpace(token)
			}))
			for _, token := range kws {
				if _, ok := keywords[token]; !ok {
					keywords[token] = struct{}{}
				}
			}
		}
		doc.KeyWords = make([]string, 0)
		for kw := range keywords {
			doc.KeyWords = append(doc.KeyWords, kw)
		}
		err := index.MapKeywordsDoc(doc.KeyWords, doc.ID)
		if err != nil {
			return err
		}
		b, err := json.Marshal(doc)
		if err != nil {
			return err
		}
		keys[i] = doc.ID
		values[i] = b
	}
	index.rwMutex.Lock()
	index.DocNum += uint64(totalBulk)
	index.UpdateAt = time.Now()
	index.rwMutex.Unlock()
	//write index to db.
	if err := index.UpdateMetadata(); err != nil {
		return err
	}
	if err := index.store.Batch(keys, values); err != nil {
		return err
	}
	return nil
}

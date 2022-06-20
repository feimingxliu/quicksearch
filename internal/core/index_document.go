package core

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/maps"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
	"strings"
	"time"
)

//TODO: support update document

func (index *Index) IndexDocument(doc *Document) error {
	if doc == nil {
		return nil
	}
	doc.IndexName = index.Name
	doc.Index = index
	//update index metadata.
	index.SetTimestamp(doc.Timestamp.UnixNano())
	index.rwMutex.Lock()
	index.DocNum++
	index.UpdateAt = time.Now()
	index.rwMutex.Unlock()
	//write to db.
	if err := index.UpdateMetadata(); err != nil {
		return err
	}
	//extract tokens and add or update inverted index.
	flatDoc := maps.Flatten(doc.Source)
	keywords := make(map[string]struct{})
	for field, value := range flatDoc {
		doc.Fields = append(doc.Fields, field)
		//this casts both string and number to string.
		s := fmt.Sprint(value)
		//filter the key char '/' and blank token.
		kws := slices.RemoveEmptyStr(slices.FilterStr(tokenizer.Tokenize(s), func(token string) string {
			return strings.Trim(strings.TrimSpace(token), "/")
		}))
		for _, token := range kws {
			if _, ok := keywords[token]; !ok {
				keywords[token] = struct{}{}
			}
		}
	}
	for kw := range keywords {
		doc.KeyWords = append(doc.KeyWords, kw)
	}
	err := InvertedIndex.MapKeywordsDoc(doc.KeyWords, doc.ID)
	if err != nil {
		return err
	}
	//write doc into db.
	b, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = db.Set(docKey(doc.IndexName, doc.ID), b)
	if err != nil {
		return err
	}
	return nil
}

func docKey(indexName string, docID string) string {
	return "/index/" + indexName + "/" + docID
}

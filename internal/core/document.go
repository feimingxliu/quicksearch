package core

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/maps"
	"strings"
	"time"
)

type Document struct {
	ID        string                 `json:"_id"`
	IndexName string                 `json:"_index"`
	Index     *Index                 `json:"-"`
	Fields    []string               `json:"fields"`
	KeyWords  []string               `json:"key_words"`
	Timestamp time.Time              `json:"@timestamp"`
	Source    map[string]interface{} `json:"_source"`
}

//TODO: support update document

func (index *Index) IndexDocument(doc *Document) error {
	if doc == nil {
		return nil
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = db.Set(docKey(doc.IndexName, doc.ID), b)
	if err != nil {
		return err
	}
	//update metadata.
	index.SetTimestamp(doc.Timestamp.UnixNano())
	index.rwMutex.Lock()
	index.DocNum++
	index.UpdateAt = time.Now()
	index.rwMutex.Unlock()
	//write to db.
	if err = index.UpdateMetadata(); err != nil {
		return err
	}
	//extract tokens and add or update inverted index.
	flatDoc := maps.Flatten(doc.Source)
	keywords := make(map[string]struct{})
	for _, value := range flatDoc {
		//this casts both string and number to string.
		s := fmt.Sprint(value)
		for _, token := range tokenizer.Tokenize(s) {
			//filter the key char '/' and blank token.
			if token = strings.Trim(token, "/"); len(token) == 0 {
				continue
			}
			if _, ok := keywords[token]; !ok {
				keywords[token] = struct{}{}
			}
		}
	}
	err = InvertedIndex.MapKeywordsDoc(keywords, doc.ID)
	if err != nil {
		return err
	}
	return nil
}

func docKey(indexName string, docID string) string {
	return "/index/" + indexName + "/" + docID
}

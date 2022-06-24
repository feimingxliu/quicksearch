package core

import (
	"context"
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/maps"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

//IndexDocument adds a doc to the index.
func (index *Index) IndexDocument(doc *Document) error {
	if doc == nil {
		return nil
	}
	if err := index.Open(); err != nil {
		return err
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
	for _, value := range flatDoc {
		//this casts both string and number to string.
		s := fmt.Sprint(value)
		//filter blank token.
		kws := slices.RemoveEmptyStr(slices.FilterStr(index.tokenizer.Tokenize(s), func(token string) string {
			return strings.TrimSpace(token)
		}))
		for _, token := range kws {
			if _, ok := keywords[token]; !ok {
				keywords[token] = struct{}{}
			}
		}
	}
	doc.KeyWords = make([]string, 0, len(keywords))
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

//UpdateDocument updates a existing doc in the index.
func (index *Index) UpdateDocument(doc *Document) error {
	if doc == nil {
		return nil
	}
	if err := index.Open(); err != nil {
		return err
	}
	//remove old keyword doc map.
	if bdoc, err := index.store.Get(doc.ID); err == nil {
		//shadowed doc.
		doc := new(Document)
		if err = json.Unmarshal(bdoc, doc); err != nil {
			return err
		}
		if err = index.UnMapKeywordsDoc(doc.KeyWords, doc.ID); err != nil {
			return err
		}
	} else {
		return errors.ErrDocumentNotFound
	}
	doc.IndexName = index.Name
	doc.Index = index
	//update index metadata.
	index.SetTimestamp(doc.Timestamp.UnixNano())
	index.rwMutex.Lock()
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
		//filter blank token.
		kws := slices.RemoveEmptyStr(slices.FilterStr(index.tokenizer.Tokenize(s), func(token string) string {
			return strings.TrimSpace(token)
		}))
		for _, token := range kws {
			if _, ok := keywords[token]; !ok {
				keywords[token] = struct{}{}
			}
		}
	}
	doc.KeyWords = make([]string, 0, len(keywords))
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

//RetrieveDocument returns the doc with docID from index.
func (index *Index) RetrieveDocument(docID string) (*Document, error) {
	if bdoc, err := index.store.Get(docID); err == nil {
		//shadowed doc.
		doc := new(Document)
		if err = json.Unmarshal(bdoc, doc); err != nil {
			return nil, err
		}
		return doc, nil
	} else {
		if err == errors.ErrKeyNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
}

//DeleteDocument deletes the doc from index.
func (index *Index) DeleteDocument(docID string) error {
	if len(docID) == 0 {
		return nil
	}
	if err := index.Open(); err != nil {
		return err
	}
	//check doc exists and update the inverted index.
	if bdoc, err := index.store.Get(docID); err == nil {
		doc := new(Document)
		if err = json.Unmarshal(bdoc, doc); err != nil {
			return err
		}
		if err = index.UnMapKeywordsDoc(doc.KeyWords, doc.ID); err != nil {
			return err
		}
	} else {
		return nil
	}
	//update index metadata.
	index.rwMutex.Lock()
	index.DocNum--
	index.UpdateAt = time.Now()
	index.rwMutex.Unlock()
	if err := index.UpdateMetadata(); err != nil {
		return err
	}
	//delete the doc.
	return index.store.Delete(docID)
}

//BulkDocuments adds docs to the index, it does not check for updated doc.
func (index *Index) BulkDocuments(docs []*Document) error {
	if len(docs) == 0 {
		return nil
	}
	if err := index.Open(); err != nil {
		return err
	}
	totalBulk := len(docs)
	g, _ := errgroup.WithContext(context.Background())
	for _, doc := range docs {
		doc.IndexName = index.Name
		doc.Index = index
		//update index metadata.
		index.SetTimestamp(doc.Timestamp.UnixNano())
		//extract tokens and add or update inverted index.
		flatDoc := maps.Flatten(doc.Source)
		keywords := make([]string, 0)
		for _, value := range flatDoc {
			//this casts both string and number to string.
			s := fmt.Sprint(value)
			//filter the key char '/' and blank token.
			kws := slices.RemoveEmptyStr(slices.FilterStr(index.tokenizer.Tokenize(s), func(token string) string {
				return strings.TrimSpace(token)
			}))
			keywords = append(keywords, kws...)
		}
		doc.KeyWords = slices.DistinctStr(keywords)
		copyDoc := doc
		//Note: use Goroutine may run out of memory when bulk large quantity of docs.
		g.Go(func() error {
			err := index.MapKeywordsDoc(copyDoc.KeyWords, copyDoc.ID)
			if err != nil {
				return err
			}
			return nil
		})
		//write doc to db.
		g.Go(func() error {
			b, err := json.Marshal(copyDoc)
			if err != nil {
				return err
			}
			if err = index.store.Batch([]string{copyDoc.ID}, [][]byte{b}); err != nil {
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	index.rwMutex.Lock()
	index.DocNum += uint64(totalBulk)
	index.UpdateAt = time.Now()
	index.rwMutex.Unlock()
	//write index to db.
	if err := index.UpdateMetadata(); err != nil {
		return err
	}
	return nil
}

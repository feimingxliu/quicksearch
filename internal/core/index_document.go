package core

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/document"
	imapping "github.com/blevesearch/bleve/v2/mapping"
	bindex "github.com/blevesearch/bleve_index_api"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"time"
)

type Document struct {
	Index       string      `json:"_index"`
	ID          string      `json:"_id"`
	Version     int64       `json:"_version"`
	SeqNo       int64       `json:"_seq_no"`
	PrimaryTerm int64       `json:"_primary_term"`
	Found       bool        `json:"found"`
	Source      interface{} `json:"_source"`
}

// IndexOrUpdateDocument indexes or update a document refers to `index`
func (index *Index) IndexOrUpdateDocument(docID string, source map[string]interface{}) error {
	shard := index.getDocShard(docID)
	doc, err := index.buildBleveDocument(docID, source, nil)
	if err != nil {
		return err
	}
	idx, err := shard.Indexer.Advanced()
	if err != nil {
		return err
	}
	return idx.Update(doc)
}

// UpdateDocumentPartially can update part fields of indexed document.
func (index *Index) UpdateDocumentPartially(docID string, fields map[string]interface{}) error {
	// check if exists
	doc, err := index.GetDocument(docID)
	if err != nil {
		return err
	}
	if !doc.Found {
		return errors.ErrDocumentNotFound
	}
	// this assert must success
	source := doc.Source.(map[string]interface{})
	for k, v := range fields {
		source[k] = v
	}
	return index.IndexOrUpdateDocument(docID, source)
}

// GetDocument returns the doc associated with docID
func (index *Index) GetDocument(docID string) (*Document, error) {
	doc := &Document{
		Index:       index.Name,
		ID:          docID,
		Version:     1,
		SeqNo:       1,
		PrimaryTerm: 1,
		Found:       false,
	}
	shard := index.getDocShard(docID)
	bdoc, err := shard.Indexer.Document(docID)
	if err != nil {
		return doc, err
	}
	if bdoc == nil {
		return doc, errors.ErrDocumentNotFound
	}
	source := make(map[string]interface{})
	bdoc.VisitFields(func(field bindex.Field) {
		if field.Name() == "_source" {
			err = json.Unmarshal(field.Value(), &source)
		}
	})
	doc.Source = source
	doc.Found = true
	return doc, err
}

// DeleteDocument try to delete the document from index, do not check if it exists
func (index *Index) DeleteDocument(docID string) error {
	shard := index.getDocShard(docID)
	return shard.Indexer.Delete(docID)
}

func (index *Index) buildBleveDocument(docID string, source map[string]interface{}, mapping imapping.IndexMapping) (*document.Document, error) {
	// add `@timestamp` field
	var err error
	doc := document.NewDocument(docID)
	if mapping != nil {
		if err = mapping.MapDocument(doc, source); err != nil {
			return nil, err
		}
		return doc, nil
	}
	if index.Mapping == nil {
		mapping = bleve.NewIndexMapping()
	} else {
		mapping, err = buildIndexMapping(index.Mapping)
		if err != nil {
			return nil, err
		}
	}
	if err = mapping.MapDocument(doc, source); err != nil {
		return nil, err
	}
	// add more fields
	doc.AddIDField()
	b, _ := json.Marshal(source)
	sf := document.NewTextFieldWithIndexingOptions("_source", nil, b, bindex.StoreField)
	doc.AddField(sf)
	dtf, err := document.NewDateTimeField("@timestamp", nil, time.Now())
	if err != nil {
		return nil, err
	}
	doc.AddField(dtf)
	cf := document.NewCompositeFieldWithIndexingOptions("_all", true, nil, []string{"_id", "_index", "_source", "@timestamp"}, bindex.IndexField)
	doc.AddField(cf)
	return doc, nil
}

func (index *Index) getDocShard(docID string) *IndexShard {
	shardID := util.BytesModInt([]byte(docID), index.NumberOfShards)
	return index.Shards[shardID]
}

package core

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/document"
	imapping "github.com/blevesearch/bleve/v2/mapping"
	"github.com/feimingxliu/quicksearch/pkg/util"
)

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

func (index *Index) buildBleveDocument(docID string, source map[string]interface{}, mapping imapping.IndexMapping) (*document.Document, error) {
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
	return doc, nil
}

func (index *Index) getDocShard(docID string) *IndexShard {
	shardID := util.BytesModInt([]byte(docID), index.NumberOfShards)
	return index.Shards[shardID]
}

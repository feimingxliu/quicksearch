package core

import (
	"github.com/blevesearch/bleve/v2"
)

const (
	DefaultNumberOfShards = 5 // default number of shards.
)

type IndexShard struct {
	ID          int         `json:"id"`           // shard's id
	DocNum      uint64      `json:"doc_num"`      // doc's number in shard
	StorageSize uint64      `json:"storage_size"` // shard file size
	Indexer     bleve.Index `json:"-"`            // a shard map to a bleve index
}

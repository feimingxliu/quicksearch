package core

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/util/uuid"
)

// BulkIndex bulk index(update if exists) docs into index
func (index *Index) BulkIndex(docs []map[string]interface{}) error {
	if len(docs) == 0 {
		return nil
	}
	batchSize := uint32(config.Global.Engine.DefaultBatchSize)
	var currentBatch uint32
	batch := make(map[int]*bleve.Batch, index.NumberOfShards)
	mapping, err := buildIndexMapping(index.Mapping)
	if err != nil {
		return err
	}
	for _, mdoc := range docs {
		docID := uuid.GetUUID()
		shard := index.getDocShard(docID)
		bleveIndex := shard.Indexer
		if batch[shard.ID] == nil {
			batch[shard.ID] = bleveIndex.NewBatch()
		}
		bdoc, err := index.buildBleveDocument(docID, mdoc, mapping)
		if err != nil {
			return err
		}
		err = batch[shard.ID].IndexAdvanced(bdoc)
		if err != nil {
			return err
		}
		currentBatch++
		// execute the batch
		if currentBatch >= batchSize {
			err = bleveIndex.Batch(batch[shard.ID])
			if err != nil {
				return err
			}
			// do not forget resetting the batch
			batch[shard.ID].Reset()
		}
	}
	// execute remaining in the batches
	for shardId, bat := range batch {
		bleveIndex := index.Shards[shardId].Indexer
		err = bleveIndex.Batch(bat)
		if err != nil {
			return err
		}
	}
	return nil
}

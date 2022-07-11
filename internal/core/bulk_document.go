package core

import (
	"bufio"
	"github.com/blevesearch/bleve/v2"
	imapping "github.com/blevesearch/bleve/v2/mapping"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/uuid"
	"io"
	"time"
)

type BulkResult struct {
	Took   time.Duration    `json:"took"`
	Errors bool             `json:"errors,omitempty"`
	Items  []BulkResultItem `json:"items"`
}

type BulkResultItem struct {
	Index  *BulkActionResult `json:"index,omitempty"`
	Delete *BulkActionResult `json:"delete,omitempty"`
	Create *BulkActionResult `json:"create,omitempty"`
	Update *BulkActionResult `json:"update,omitempty"`
}

type BulkActionResult struct {
	Index       string      `json:"_index"`
	ID          string      `json:"_id"`
	Version     int64       `json:"_version"`
	Result      string      `json:"result,omitempty"` // 'created', 'deleted', 'indexed', 'updated'
	Shards      *BulkShards `json:"_shards"`
	Status      int64       `json:"status"` // 200 => indexed, updated, deleted	201 => created	404 => doc not found
	Error       interface{} `json:"error,omitempty"`
	SeqNo       int64       `json:"_seq_no"`
	PrimaryTerm int64       `json:"_primary_term"`
}

type BulkShards struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Failed     int64 `json:"failed"`
}

type BulkAction struct {
	Index  *BulkActionDetail `json:"index,omitempty"`
	Create *BulkActionDetail `json:"create,omitempty"`
	Update *BulkActionDetail `json:"update,omitempty"`
	Delete *BulkActionDetail `json:"delete,omitempty"`
}

type BulkActionDetail struct {
	Index string `json:"_index"`
	ID    string `json:"_id"`
}

// Bulk reads data from reader and execute bulk actions defined in reader.
// The first line looks like {"$action":{"_index": "$index", "_id": "$docID"}},
// the $action can be `create`, `delete`, `index`, `update`. Note that here's
// update don't support update document partially  because of performance.
// If $index is empty, the targetIndex will be used.
func Bulk(targetIndex string, reader io.Reader) (*BulkResult, error) {
	var (
		startTime        = time.Now()
		err              error
		scanner          = bufio.NewScanner(reader)
		bulkResult       = &BulkResult{}
		mapping          = make(map[string]imapping.IndexMapping) // index => IndexMapping
		batch            = make(map[bleve.Index]*bleve.Batch)     // index => Batch
		batchSize        = uint32(config.Global.Engine.DefaultBatchSize)
		currentBatchSize uint32
		action           = new(BulkAction)
		nextLineIsData   bool
		indexName        string
		index            *Index
		bindex           bleve.Index
		data             = make(map[string]interface{})
	)

	defer func() {
		bulkResult.Took = time.Since(startTime)
		if err != nil {
			bulkResult.Errors = true
		}
	}()

	for scanner.Scan() {
		if !nextLineIsData {
			err = json.Unmarshal(scanner.Bytes(), action)
			if err != nil {
				return bulkResult, errors.ErrBulkDataFormat
			}
			if action.Delete != nil {
				indexName = action.Delete.Index
				if indexName == "" {
					indexName = targetIndex
					if indexName == "" {
						return bulkResult, errors.ErrBulkDataFormat
					}
				}
				index, err := GetIndex(indexName)
				if err != nil {
					if err == errors.ErrIndexNotFound {
						bulkResult.Items = append(bulkResult.Items, BulkResultItem{
							Delete: NewBulkActionResult(indexName, action.Delete.ID, "index_not_found", 404, err, int64(len(bulkResult.Items))),
						})
					} else {
						return bulkResult, err
					}
					continue
				}
				docID := action.Delete.ID
				bindex = index.getDocShard(docID).Indexer
				if batch[bindex] == nil {
					batch[bindex] = bindex.NewBatch()
				}
				batch[bindex].Delete(docID)
				bulkResult.Items = append(bulkResult.Items, BulkResultItem{
					Delete: NewBulkActionResult(indexName, action.Delete.ID, "deleted", 200, nil, int64(len(bulkResult.Items))),
				})
				nextLineIsData = false
				continue
			}
			if action.Index != nil || action.Create != nil || action.Update != nil {
				nextLineIsData = true
			}
		} else {
			nextLineIsData = false

			err = json.Unmarshal(scanner.Bytes(), &data)
			if err != nil {
				return bulkResult, errors.ErrBulkDataFormat
			}

			if action.Index != nil || action.Create != nil || action.Update != nil {
				var docID string
				switch {
				case action.Index != nil:
					indexName = action.Index.Index
					docID = action.Index.ID
					if indexName == "" {
						indexName = targetIndex
						if indexName == "" {
							return bulkResult, errors.ErrBulkDataFormat
						}
					}
					if docID == "" {
						docID = uuid.GetUUID()
					}
					bulkResult.Items = append(bulkResult.Items, BulkResultItem{
						Index: NewBulkActionResult(indexName, docID, "indexed", 200, nil, int64(len(bulkResult.Items))),
					})
				case action.Create != nil:
					indexName = action.Create.Index
					docID = action.Create.ID
					if indexName == "" {
						indexName = targetIndex
						if indexName == "" {
							return bulkResult, errors.ErrBulkDataFormat
						}
					}
					if docID == "" {
						docID = uuid.GetUUID()
					}
					bulkResult.Items = append(bulkResult.Items, BulkResultItem{
						Create: NewBulkActionResult(indexName, docID, "created", 201, nil, int64(len(bulkResult.Items))),
					})
				case action.Update != nil:
					indexName = action.Update.Index
					docID = action.Update.ID
					if indexName == "" {
						indexName = targetIndex
						if indexName == "" {
							return bulkResult, errors.ErrBulkDataFormat
						}
					}
					if docID == "" {
						docID = uuid.GetUUID()
					}
					bulkResult.Items = append(bulkResult.Items, BulkResultItem{
						Update: NewBulkActionResult(indexName, docID, "updated", 200, nil, int64(len(bulkResult.Items))),
					})
				}
				index, err = NewIndex(WithName(indexName))
				if err != nil {
					return bulkResult, err
				}
				bindex = index.getDocShard(docID).Indexer
				if batch[bindex] == nil {
					batch[bindex] = bindex.NewBatch()
				}
				if mapping[indexName] == nil {
					mp, err := buildIndexMapping(index.Mapping)
					if err != nil {
						return bulkResult, err
					}
					mapping[indexName] = mp
				}
				bdoc, err := index.buildBleveDocument(docID, data, mapping[indexName])
				if err != nil {
					return bulkResult, err
				}
				err = batch[bindex].IndexAdvanced(bdoc)
				if err != nil {
					return bulkResult, err
				}

				currentBatchSize++
				if currentBatchSize >= batchSize {
					for bindex, bat := range batch {
						err = bindex.Batch(bat)
						if err != nil {
							return bulkResult, err
						}
						// do not forget to reset batch
						bat.Reset()
					}
					currentBatchSize = 0
				}
				continue
			}

			action = new(BulkAction)
		}
	}

	// bulk the remaining
	for bindex, bat := range batch {
		if bat.Size() > 0 {
			err = bindex.Batch(bat)
			if err != nil {
				return bulkResult, err
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return bulkResult, err
	}

	return bulkResult, err
}

func NewBulkActionResult(index string, docID string, result string, status int64, err interface{}, seqNo int64) *BulkActionResult {
	return &BulkActionResult{
		Index:       index,
		ID:          docID,
		Version:     1,
		Result:      result,
		Status:      status,
		Error:       err,
		SeqNo:       seqNo,
		PrimaryTerm: 1,
		Shards: &BulkShards{
			Total:      1,
			Successful: 1,
			Failed:     0,
		},
	}
}

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
			// do not forget to reset the batch
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

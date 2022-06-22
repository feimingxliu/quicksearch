package core

import (
	"github.com/feimingxliu/quicksearch/pkg/util/uuid"
	"time"
)

type Document struct {
	ID        string                 `json:"_id"`
	IndexName string                 `json:"_index"`
	Index     *Index                 `json:"-"`
	KeyWords  []string               `json:"key_words"`
	Timestamp time.Time              `json:"@timestamp"`
	Source    map[string]interface{} `json:"_source"`
}

func NewDocument(source map[string]interface{}) *Document {
	return &Document{
		ID:        uuid.GetXID(),
		Timestamp: time.Now(),
		Source:    source,
	}
}

func (doc *Document) WithID(id string) *Document {
	doc.ID = id
	return doc
}

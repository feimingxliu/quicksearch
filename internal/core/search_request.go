package core

import (
	"github.com/blevesearch/bleve/v2/search/query"
	"time"
)

type SearchRequest struct {
	Query            query.Query              `json:"query"`
	Size             int                      `json:"size"`
	From             int                      `json:"from"`
	Highlight        bool                     `json:"highlight"`
	Fields           []string                 `json:"fields"`
	Facets           map[string]*FacetRequest `json:"facets"`
	Explain          bool                     `json:"explain"`
	Sort             []string                 `json:"sort"`
	IncludeLocations bool                     `json:"includeLocations"`
	SearchAfter      []string                 `json:"search_after"`
	SearchBefore     []string                 `json:"search_before"`
}

type FacetRequest struct {
	Size           int             `json:"size"`
	Field          string          `json:"field"`
	NumericRanges  []NumericRange  `json:"numeric_ranges,omitempty"`
	DateTimeRanges []DateTimeRange `json:"date_ranges,omitempty"`
}

type NumericRange struct {
	Name string  `json:"name,omitempty"`
	Min  float64 `json:"min,omitempty"`
	Max  float64 `json:"max,omitempty"`
}

type DateTimeRange struct {
	Name  string    `json:"name,omitempty"`
	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`
}

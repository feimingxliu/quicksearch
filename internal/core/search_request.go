package core

import (
	stdjson "encoding/json"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
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

func (r *SearchRequest) UnmarshalJSON(input []byte) error {
	var temp struct {
		Q                stdjson.RawMessage       `json:"query"`
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
	err := json.Unmarshal(input, &temp)
	if err != nil {
		return err
	}
	if temp.Size <= 0 {
		r.Size = config.Global.Engine.DefaultSearchResultSize
	} else {
		r.Size = temp.Size
	}
	r.From = temp.From
	r.Explain = temp.Explain
	r.Highlight = temp.Highlight
	r.Fields = temp.Fields
	r.Facets = temp.Facets
	r.IncludeLocations = temp.IncludeLocations
	r.SearchAfter = temp.SearchAfter
	r.SearchBefore = temp.SearchBefore
	r.Sort = temp.Sort
	r.Query, err = query.ParseQuery(temp.Q)
	if err != nil {
		return err
	}
	if r.From < 0 {
		r.From = 0
	}
	return nil
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

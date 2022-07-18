package core

import (
	"github.com/blevesearch/bleve/v2/search"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"math"
	"time"
)

type SearchResult struct {
	Status    Status                  `json:"status"`
	Request   *SearchRequest          `json:"request"`
	Hits      Hits                    `json:"hits"`
	TotalHits uint64                  `json:"total_hits"`
	MaxScore  float64                 `json:"max_score"`
	Took      time.Duration           `json:"took"`
	Facets    map[string]*FacetResult `json:"facets,omitempty"`
}

type Status struct {
	Total      int `json:"total"`
	Failed     int `json:"failed"`
	Successful int `json:"successful"`
}

type FacetResult struct {
	Field         string              `json:"field"`
	Total         int                 `json:"total"`
	Missing       int                 `json:"missing"`
	Other         int                 `json:"other"`
	DateRanges    []DateRangeFacet    `json:"date_ranges,omitempty"`
	NumericRanges []NumericRangeFacet `json:"numeric_ranges,omitempty"`
	Terms         []TermFacet         `json:"terms,omitempty"`
}

type DateRangeFacet struct {
	Name  string  `json:"name"`
	Start *string `json:"start,omitempty"`
	End   *string `json:"end,omitempty"`
	Count int     `json:"count"`
}

type NumericRangeFacet struct {
	Name  string   `json:"name"`
	Min   *float64 `json:"min"`
	Max   *float64 `json:"max"`
	Count int      `json:"count"`
}

type TermFacet struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

type f64 float64

type Hit struct {
	Index       string                      `json:"_index"`
	ID          string                      `json:"_id"`
	Score       f64                         `json:"_score"`
	Sort        []string                    `json:"_sort"`
	Timestamp   string                      `json:"@timestamp"`
	Explanation *search.Explanation         `json:"_explanation,omitempty"`
	Locations   search.FieldTermLocationMap `json:"_locations,omitempty"`
	Fragments   search.FieldFragmentMap     `json:"_fragments,omitempty"`
	Source      map[string]interface{}      `json:"_source,omitempty"`
	Fields      map[string]interface{}      `json:"_fields,omitempty"`
}

func (f f64) MarshalJSON() ([]byte, error) {
	if !math.IsNaN(float64(f)) {
		return json.Marshal(float64(f))
	} else {
		return json.Marshal(0)
	}
}

type Hits []*Hit

func (h Hits) Len() int {
	return len(h)
}

func (h Hits) Less(i, j int) bool {
	return h[i].Score > h[j].Score
}

func (h Hits) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

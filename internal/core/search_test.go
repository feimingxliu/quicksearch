package core

import (
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"log"
	"testing"
)

func TestSearchQueryStringQuery(t *testing.T) {
	prepare(t)
	defer clean(t)
	bulkDocument(t, 10000, false)
	index, err := GetIndex(indexName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		log.Println("Delete Index.")
		if err := index.Delete(); err != nil {
			t.Fatal(err)
		}
	}()
	searchRequest := &SearchRequest{
		Query:     &QueryStringQuery{Query: "数学，是研究数量、结构以及空间等概念及其变化的一门学科，从某种角度看属於形式科学的一种"},
		Size:      1,
		From:      0,
		Highlight: false,
		Fields:    []string{"*"},
		Facets: map[string]*FacetRequest{
			"terms": {
				Size:  10,
				Field: "id",
			},
		},
		Explain:          false,
		IncludeLocations: false,
	}
	res, err := index.Search(searchRequest)
	if err != nil {
		t.Fatal(err)
	}
	json.Print("result", res)
}

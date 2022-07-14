package core

import (
	imapping "github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/blevesearch/bleve/v2/search/query"
	bindex "github.com/blevesearch/bleve_index_api"
)

type QueryStringQuery struct {
	Query string  `json:"query"`
	Boost float64 `json:"boost"`
}

func (q *QueryStringQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(q.Boost)
	qq := query.QueryStringQuery{
		Query:    q.Query,
		BoostVal: &boost,
	}
	return qq.Searcher(i, m, options)
}

// todo: support more query

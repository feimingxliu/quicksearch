package core

import (
	imapping "github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/blevesearch/bleve/v2/search/query"
	bindex "github.com/blevesearch/bleve_index_api"
	"strings"
	"time"
)

type QueryStringQuery struct {
	Query string  `json:"query"`
	Boost float64 `json:"boost,omitempty"`
}

func (q *QueryStringQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(q.Boost)
	qq := query.QueryStringQuery{
		Query:    q.Query,
		BoostVal: &boost,
	}
	return qq.Searcher(i, m, options)
}

type TermQuery struct {
	Term     string  `json:"term"`
	FieldVal string  `json:"field,omitempty"`
	BoostVal float64 `json:"boost,omitempty"`
}

func (t *TermQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(t.BoostVal)
	tq := query.TermQuery{
		Term:     t.Term,
		FieldVal: t.FieldVal,
		BoostVal: &boost,
	}
	return tq.Searcher(i, m, options)
}

type MatchQuery struct {
	Match     string  `json:"match"`
	FieldVal  string  `json:"field,omitempty"`
	Analyzer  string  `json:"analyzer,omitempty"`
	BoostVal  float64 `json:"boost,omitempty"`
	Prefix    int     `json:"prefix_length"`
	Fuzziness int     `json:"fuzziness"`
	Operator  string  `json:"operator,omitempty"`
}

func (mq *MatchQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(mq.BoostVal)
	q := query.MatchQuery{
		Match:     mq.Match,
		FieldVal:  mq.FieldVal,
		Analyzer:  mq.Analyzer,
		BoostVal:  &boost,
		Prefix:    mq.Prefix,
		Fuzziness: mq.Fuzziness,
		Operator: func() query.MatchQueryOperator {
			switch strings.ToLower(mq.Operator) {
			case "and":
				return query.MatchQueryOperatorAnd
			case "or":
				return query.MatchQueryOperatorOr
			default:
				return query.MatchQueryOperatorOr
			}
		}(),
	}
	return q.Searcher(i, m, options)
}

type PhraseQuery struct {
	Terms    []string `json:"terms"`
	Field    string   `json:"field,omitempty"`
	BoostVal float64  `json:"boost,omitempty"`
}

func (p *PhraseQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(p.BoostVal)
	pq := query.PhraseQuery{
		Terms:    p.Terms,
		Field:    p.Field,
		BoostVal: &boost,
	}
	return pq.Searcher(i, m, options)
}

type MatchPhraseQuery struct {
	MatchPhrase string  `json:"match_phrase"`
	FieldVal    string  `json:"field,omitempty"`
	Analyzer    string  `json:"analyzer,omitempty"`
	BoostVal    float64 `json:"boost,omitempty"`
}

func (mp *MatchPhraseQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(mp.BoostVal)
	mpq := query.MatchPhraseQuery{
		MatchPhrase: mp.MatchPhrase,
		FieldVal:    mp.FieldVal,
		Analyzer:    mp.Analyzer,
		BoostVal:    &boost,
	}
	return mpq.Searcher(i, m, options)
}

type FuzzyQuery struct {
	Term      string  `json:"term"`
	Prefix    int     `json:"prefix_length"`
	Fuzziness int     `json:"fuzziness"`
	FieldVal  string  `json:"field,omitempty"`
	BoostVal  float64 `json:"boost,omitempty"`
}

func (f *FuzzyQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(f.BoostVal)
	fq := query.FuzzyQuery{
		Term:      f.Term,
		Prefix:    f.Prefix,
		Fuzziness: f.Fuzziness,
		FieldVal:  f.FieldVal,
		BoostVal:  &boost,
	}
	return fq.Searcher(i, m, options)
}

type ConjunctionQuery struct {
	Conjuncts []query.Query `json:"conjuncts"`
	BoostVal  float64       `json:"boost,omitempty"`
}

func (c *ConjunctionQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(c.BoostVal)
	cq := query.ConjunctionQuery{
		Conjuncts: c.Conjuncts,
		BoostVal:  &boost,
	}
	return cq.Searcher(i, m, options)
}

type DisjunctionQuery struct {
	Disjuncts []query.Query `json:"disjuncts"`
	BoostVal  float64       `json:"boost,omitempty"`
	Min       float64       `json:"min"`
}

func (d *DisjunctionQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(d.BoostVal)
	dq := query.DisjunctionQuery{
		Disjuncts: d.Disjuncts,
		BoostVal:  &boost,
		Min:       d.Min,
	}
	return dq.Searcher(i, m, options)
}

type BooleanQuery struct {
	Must     query.Query `json:"must,omitempty"`
	Should   query.Query `json:"should,omitempty"`
	MustNot  query.Query `json:"must_not,omitempty"`
	BoostVal float64     `json:"boost,omitempty"`
}

func (b *BooleanQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(b.BoostVal)
	bq := query.BooleanQuery{
		Must:     b.Must,
		Should:   b.Should,
		MustNot:  b.MustNot,
		BoostVal: &boost,
	}
	return bq.Searcher(i, m, options)
}

type NumericRangeQuery struct {
	Min          *float64 `json:"min,omitempty"`
	Max          *float64 `json:"max,omitempty"`
	InclusiveMin *bool    `json:"inclusive_min,omitempty"`
	InclusiveMax *bool    `json:"inclusive_max,omitempty"`
	FieldVal     string   `json:"field,omitempty"`
	BoostVal     float64  `json:"boost,omitempty"`
}

func (n *NumericRangeQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(n.BoostVal)
	nrq := query.NumericRangeQuery{
		Min:          n.Min,
		Max:          n.Max,
		InclusiveMin: n.InclusiveMin,
		InclusiveMax: n.InclusiveMax,
		FieldVal:     n.FieldVal,
		BoostVal:     &boost,
	}
	return nrq.Searcher(i, m, options)
}

type DateRangeQuery struct {
	Start          time.Time `json:"start,omitempty"`
	End            time.Time `json:"end,omitempty"`
	InclusiveStart *bool     `json:"inclusive_start,omitempty"`
	InclusiveEnd   *bool     `json:"inclusive_end,omitempty"`
	FieldVal       string    `json:"field,omitempty"`
	BoostVal       float64   `json:"boost,omitempty"`
}

func (d *DateRangeQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(d.BoostVal)
	drq := query.DateRangeQuery{
		Start:          query.BleveQueryTime{d.Start},
		End:            query.BleveQueryTime{d.End},
		InclusiveStart: d.InclusiveStart,
		InclusiveEnd:   d.InclusiveEnd,
		FieldVal:       d.FieldVal,
		BoostVal:       &boost,
	}
	return drq.Searcher(i, m, options)
}

type MatchAllQuery struct {
	BoostVal float64 `json:"boost,omitempty"`
}

func (mq *MatchAllQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(mq.BoostVal)
	maq := query.MatchAllQuery{BoostVal: &boost}
	return maq.Searcher(i, m, options)
}

type MatchNoneQuery struct {
	BoostVal float64 `json:"boost,omitempty"`
}

func (mq *MatchNoneQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(mq.BoostVal)
	mnq := query.MatchNoneQuery{BoostVal: &boost}
	return mnq.Searcher(i, m, options)
}

type DocIDQuery struct {
	IDs      []string `json:"ids"`
	BoostVal float64  `json:"boost,omitempty"`
}

func (d *DocIDQuery) Searcher(i bindex.IndexReader, m imapping.IndexMapping, options search.SearcherOptions) (search.Searcher, error) {
	boost := query.Boost(d.BoostVal)
	dq := query.DocIDQuery{
		IDs:      d.IDs,
		BoostVal: &boost,
	}
	return dq.Searcher(i, m, options)
}

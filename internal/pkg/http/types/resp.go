package types

import "github.com/feimingxliu/quicksearch/internal/core"

type Common struct {
	Acknowledged bool   `json:"_acknowledged"`
	Error        string `json:"_error,omitempty"`
}

type Index struct {
	Common
	*core.Index
}

type Indices struct {
	Common
	Indices []*core.Index `json:"indices"`
}

type IndexDocument struct {
	Common
	Index string `json:"_index"`
	ID    string `json:"_id"`
	Type  string `json:"_type"`
}

type SearchResult struct {
	Common
	*core.SearchResult
}

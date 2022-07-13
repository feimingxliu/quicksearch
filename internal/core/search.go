package core

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
)

// Search performs search in specified index.
func (index *Index) Search(req *SearchRequest) (*SearchResult, error) {
	request := &bleve.SearchRequest{
		Query:            req.Query,
		Size:             req.Size,
		From:             req.From,
		Fields:           []string{"*"},
		Explain:          req.Explain,
		Sort:             search.SortOrder{&search.SortScore{Desc: true}},
		IncludeLocations: req.IncludeLocations,
		SearchAfter: func() []string {
			if len(req.SearchAfter) > 0 {
				return req.SearchAfter
			}
			return nil
		}(),
		SearchBefore: func() []string {
			if len(req.SearchBefore) > 0 {
				return req.SearchBefore
			}
			return nil
		}(),
	}
	source := true
	fields := make(map[string]bool, len(request.Fields))
	if len(request.Fields) > 0 {
		request.Fields = req.Fields
		if !slices.ContainsStr(request.Fields, "*") && !slices.ContainsStr(request.Fields, "_all") {
			source = false
			for _, f := range request.Fields {
				fields[f] = true
			}
		}
	}
	if req.Size == 0 {
		request.Size = config.Global.Engine.DefaultSearchResultSize
	}
	if req.Highlight {
		request.Highlight = bleve.NewHighlight()
	}
	if len(req.Facets) > 0 {
		facets := make(map[string]*bleve.FacetRequest, len(req.Facets))
		for name, fr := range req.Facets {
			facet := bleve.NewFacetRequest(fr.Field, fr.Size)
			for _, nr := range fr.NumericRanges {
				facet.AddNumericRange(nr.Name, &nr.Min, &nr.Max)
			}
			for _, dr := range fr.DateTimeRanges {
				facet.AddDateTimeRange(dr.Name, dr.Start, dr.End)
			}
			facets[name] = facet
		}
		request.Facets = facets
	}
	if len(req.Sort) > 0 {
		so := search.ParseSortOrderStrings(req.Sort)
		request.Sort = so
	}
	indexes := make([]bleve.Index, 0, index.NumberOfShards)
	for _, shard := range index.Shards {
		indexes = append(indexes, shard.Indexer)
	}
	indexAlias := bleve.NewIndexAlias(indexes...)
	searchResult, err := indexAlias.Search(request)
	if err != nil {
		return nil, err
	}
	result := &SearchResult{
		Status: Status{
			Total:      searchResult.Status.Total,
			Failed:     searchResult.Status.Failed,
			Successful: searchResult.Status.Successful,
		},
		Request:   req,
		Hits:      make([]*Hit, 0, len(searchResult.Hits)),
		TotalHits: searchResult.Total,
		MaxScore:  searchResult.MaxScore,
		Took:      searchResult.Took,
	}
	for _, dm := range searchResult.Hits {
		hit := &Hit{
			Index:       dm.Index,
			ID:          dm.ID,
			Score:       dm.Score,
			Sort:        dm.Sort,
			Explanation: dm.Expl,
			Locations:   dm.Locations,
			Fragments:   dm.Fragments,
		}
		if source {
			hit.Source = map[string]interface{}{}
		} else {
			hit.Fields = map[string]interface{}{}
		}
		for k, v := range dm.Fields {
			switch k {
			case "@timestamp":
				if t, ok := v.(string); ok {
					hit.Timestamp = t
				}
			case "_source":
				if source {
					if s, ok := v.(string); ok {
						if err = json.Unmarshal([]byte(s), &hit.Source); err != nil {
							return nil, err
						}
					}
				}
			default:
				if !source && fields[k] {
					hit.Fields[k] = v
				}
			}
		}
		result.Hits = append(result.Hits, hit)
	}
	if searchResult.Facets != nil {
		facets := make(map[string]*FacetResult, len(searchResult.Facets))
		for name, fr := range searchResult.Facets {
			facet := &FacetResult{
				Field:   fr.Field,
				Total:   fr.Total,
				Missing: fr.Missing,
				Other:   fr.Other,
			}
			if fr.DateRanges != nil {
				facet.DateRanges = make([]DateRangeFacet, 0, len(fr.DateRanges))
				for _, drf := range fr.DateRanges {
					facet.DateRanges = append(facet.DateRanges, DateRangeFacet{
						Name:  drf.Name,
						Start: drf.Start,
						End:   drf.End,
						Count: drf.Count,
					})
				}
			}
			if fr.NumericRanges != nil {
				facet.NumericRanges = make([]NumericRangeFacet, 0, len(fr.NumericRanges))
				for _, nr := range fr.NumericRanges {
					facet.NumericRanges = append(facet.NumericRanges, NumericRangeFacet{
						Name:  nr.Name,
						Min:   nr.Min,
						Max:   nr.Max,
						Count: nr.Count,
					})
				}
			}
			if fr.Terms != nil {
				terms := fr.Terms.Terms()
				facet.Terms = make([]TermFacet, 0, len(terms))
				for _, term := range terms {
					facet.Terms = append(facet.Terms, TermFacet{
						Term:  term.Term,
						Count: term.Count,
					})
				}
			}
			facets[name] = facet
		}
		result.Facets = facets
	}

	return result, nil
}

// Search performs search in all existing indices.
func Search(req *SearchRequest) (*SearchResult, error) {
	request := &bleve.SearchRequest{
		Query:            req.Query,
		Size:             req.Size,
		From:             req.From,
		Fields:           []string{"*"},
		Explain:          req.Explain,
		Sort:             search.SortOrder{&search.SortScore{Desc: true}},
		IncludeLocations: req.IncludeLocations,
		SearchAfter: func() []string {
			if len(req.SearchAfter) > 0 {
				return req.SearchAfter
			}
			return nil
		}(),
		SearchBefore: func() []string {
			if len(req.SearchBefore) > 0 {
				return req.SearchBefore
			}
			return nil
		}(),
	}
	source := true
	fields := make(map[string]bool, len(request.Fields))
	if len(request.Fields) > 0 {
		request.Fields = req.Fields
		if !slices.ContainsStr(request.Fields, "*") && !slices.ContainsStr(request.Fields, "_all") {
			source = false
			for _, f := range request.Fields {
				fields[f] = true
			}
		}
	}
	if req.Size == 0 {
		request.Size = config.Global.Engine.DefaultSearchResultSize
	}
	if req.Highlight {
		request.Highlight = bleve.NewHighlight()
	}
	if len(req.Facets) > 0 {
		facets := make(map[string]*bleve.FacetRequest, len(req.Facets))
		for name, fr := range req.Facets {
			facet := bleve.NewFacetRequest(fr.Field, fr.Size)
			for _, nr := range fr.NumericRanges {
				facet.AddNumericRange(nr.Name, &nr.Min, &nr.Max)
			}
			for _, dr := range fr.DateTimeRanges {
				facet.AddDateTimeRange(dr.Name, dr.Start, dr.End)
			}
			facets[name] = facet
		}
		request.Facets = facets
	}
	if len(req.Sort) > 0 {
		so := search.ParseSortOrderStrings(req.Sort)
		request.Sort = so
	}
	// to get all indices(some may close), fetch all from meta db.
	indices, err := ListIndices()
	if err != nil {
		return nil, err
	}
	indexes := make([]bleve.Index, 0)
	for _, index := range indices {
		// this will search in cache first, some index may already open and exist in the engine cache.
		// the returned index is opened.
		if index, err = GetIndex(index.Name); err != nil {
			return nil, err
		}
		for _, shard := range index.Shards {
			indexes = append(indexes, shard.Indexer)
		}
	}
	if len(indexes) == 0 {
		return nil, errors.ErrIndexNotFound
	}
	indexAlias := bleve.NewIndexAlias(indexes...)
	searchResult, err := indexAlias.Search(request)
	if err != nil {
		return nil, err
	}
	result := &SearchResult{
		Status: Status{
			Total:      searchResult.Status.Total,
			Failed:     searchResult.Status.Failed,
			Successful: searchResult.Status.Successful,
		},
		Request:   req,
		Hits:      make([]*Hit, 0, len(searchResult.Hits)),
		TotalHits: searchResult.Total,
		MaxScore:  searchResult.MaxScore,
		Took:      searchResult.Took,
	}
	for _, dm := range searchResult.Hits {
		hit := &Hit{
			Index:       dm.Index,
			ID:          dm.ID,
			Score:       dm.Score,
			Sort:        dm.Sort,
			Explanation: dm.Expl,
			Locations:   dm.Locations,
			Fragments:   dm.Fragments,
		}
		if source {
			hit.Source = map[string]interface{}{}
		} else {
			hit.Fields = map[string]interface{}{}
		}
		for k, v := range dm.Fields {
			switch k {
			case "@timestamp":
				if t, ok := v.(string); ok {
					hit.Timestamp = t
				}
			case "_source":
				if source {
					if s, ok := v.(string); ok {
						if err = json.Unmarshal([]byte(s), &hit.Source); err != nil {
							return nil, err
						}
					}
				}
			default:
				if !source && fields[k] {
					hit.Fields[k] = v
				}
			}
		}
		result.Hits = append(result.Hits, hit)
	}
	if searchResult.Facets != nil {
		facets := make(map[string]*FacetResult, len(searchResult.Facets))
		for name, fr := range searchResult.Facets {
			facet := &FacetResult{
				Field:   fr.Field,
				Total:   fr.Total,
				Missing: fr.Missing,
				Other:   fr.Other,
			}
			if fr.DateRanges != nil {
				facet.DateRanges = make([]DateRangeFacet, 0, len(fr.DateRanges))
				for _, drf := range fr.DateRanges {
					facet.DateRanges = append(facet.DateRanges, DateRangeFacet{
						Name:  drf.Name,
						Start: drf.Start,
						End:   drf.End,
						Count: drf.Count,
					})
				}
			}
			if fr.NumericRanges != nil {
				facet.NumericRanges = make([]NumericRangeFacet, 0, len(fr.NumericRanges))
				for _, nr := range fr.NumericRanges {
					facet.NumericRanges = append(facet.NumericRanges, NumericRangeFacet{
						Name:  nr.Name,
						Min:   nr.Min,
						Max:   nr.Max,
						Count: nr.Count,
					})
				}
			}
			if fr.Terms != nil {
				terms := fr.Terms.Terms()
				facet.Terms = make([]TermFacet, 0, len(terms))
				for _, term := range terms {
					facet.Terms = append(facet.Terms, TermFacet{
						Term:  term.Term,
						Count: term.Count,
					})
				}
			}
			facets[name] = facet
		}
		result.Facets = facets
	}

	return result, nil
}

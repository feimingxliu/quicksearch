package core

import "time"

type SearchResult struct {
	Took     string  `json:"took"`
	TimedOut bool    `json:"timed_out"`
	MaxScore float64 `json:"max_score"`
	Hits     Hits    `json:"hits"`
	Error    string  `json:"error,omitempty"`
}

type Hits struct {
	Total Total  `json:"total"`
	Hits  []*Hit `json:"hits"`
}

type Total struct {
	Value int `json:"value"`
}

type Hit struct {
	Index     string                 `json:"_index"`
	Type      string                 `json:"_type"`
	ID        string                 `json:"_id"`
	Score     float64                `json:"_score"`
	Timestamp time.Time              `json:"@timestamp"`
	Source    map[string]interface{} `json:"_source"`
}

type HitSlice []*Hit

func (h HitSlice) Len() int {
	return len(h)
}

func (h HitSlice) Less(i, j int) bool {
	return h[i].Score > h[j].Score
}

func (h HitSlice) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

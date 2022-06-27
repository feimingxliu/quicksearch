package types

type Search struct {
	Query   string `json:"query"`
	Timeout int    `json:"timeout"`
	TopN    int    `json:"top_n"`
}

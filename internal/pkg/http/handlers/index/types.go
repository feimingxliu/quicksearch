package index

import "github.com/feimingxliu/quicksearch/internal/core"

type CreateIndex struct {
	Settings *Settings          `json:"settings"`
	Mappings *core.IndexMapping `json:"mappings"`
}

type Settings struct {
	NumberOfShards int `json:"number_of_shards"`
}

package gse

import (
	"errors"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

func analyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	tokenizerName, ok := config["tokenizer"].(string)
	if !ok {
		return nil, errors.New("must have tokenizer")
	}

	tokenizer, err := cache.TokenizerNamed(tokenizerName)
	if err != nil {
		return nil, err
	}

	az := &analysis.Analyzer{Tokenizer: tokenizer}
	return az, nil
}

func init() {
	registry.RegisterAnalyzer("gse", analyzerConstructor)
}

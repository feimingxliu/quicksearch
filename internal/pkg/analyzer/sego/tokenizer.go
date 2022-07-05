package sego

import (
	"errors"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/huichen/sego"
	"path"
)

const Name = "sego"

type SegoTokenizer struct {
	seg *sego.Segmenter
}

func NewSegoTokenizer(dictPath string) *SegoTokenizer {
	if dictPath == "" {
		dictPath = path.Join(DictDir, "dictionary.txt")
	}
	seg := new(sego.Segmenter)
	seg.LoadDictionary(dictPath)
	return &SegoTokenizer{seg: seg}
}

func (s SegoTokenizer) Tokenize(text []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	pos := 1
	segments := s.seg.Segment(text)
	for _, segment := range segments {
		typ := analysis.Ideographic
		alphaNumeric := true
		for _, r := range segment.Token().Text() {
			if r < 32 || r > 126 {
				alphaNumeric = false
				break
			}
		}
		if alphaNumeric {
			typ = analysis.AlphaNumeric
		}
		token := analysis.Token{
			Term:     []byte(segment.Token().Text()),
			Start:    segment.Start(),
			End:      segment.End(),
			Position: pos,
			Type:     typ,
		}
		result = append(result, &token)
		pos++
	}
	return result
}

func tokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	dictpath, ok := config["dict_path"].(string)
	if !ok {
		return nil, errors.New("config dict_path not found")
	}
	return NewSegoTokenizer(dictpath), nil
}

func init() {
	registry.RegisterTokenizer(Name, tokenizerConstructor)
}

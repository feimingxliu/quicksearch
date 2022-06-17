package tokenizer

import "github.com/feimingxliu/quicksearch/pkg/tokenizer/jieba"

func NewTokenizer(typ string) Tokenizer {
	switch typ {
	case "jieba":
		return jieba.NewJieBa()
	default:
		return jieba.NewJieBa()
	}
}

type Tokenizer interface {
	Tokenize(s string) []string
	Keywords(s string, topK int) []string
	KeywordsWeight(s string, topK int) []WordWeight
	Close()
}

type WordWeight struct {
	Word   string  `json:"word"`
	Weight float64 `json:"weight"`
}

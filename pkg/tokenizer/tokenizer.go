package tokenizer

import "github.com/feimingxliu/quicksearch/pkg/tokenizer/jieba"

type TokenizeType int

//types of tokenizer
const (
	Jieba TokenizeType = iota

	Default = Jieba
)

func NewTokenizer(typ TokenizeType) Tokenizer {
	switch typ {
	case Jieba:
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

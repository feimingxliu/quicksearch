package core

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	ptokenizer "github.com/feimingxliu/quicksearch/internal/pkg/tokenizer"
	"strings"
)

var tokenizer ptokenizer.Tokenizer

func init() {
	switch strings.ToLower(config.Global.Tokenizer.Type) {
	case "jieba":
		tokenizer = ptokenizer.NewTokenizer(ptokenizer.Jieba)
	default:
		tokenizer = ptokenizer.NewTokenizer(ptokenizer.Default)
	}
}

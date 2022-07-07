package analyzer

import (
	_ "github.com/feimingxliu/quicksearch/internal/pkg/analyzer/gse"
	_ "github.com/feimingxliu/quicksearch/internal/pkg/analyzer/jieba"
	_ "github.com/feimingxliu/quicksearch/internal/pkg/analyzer/sego"
)

// this package just import the underlying analyzer for init()

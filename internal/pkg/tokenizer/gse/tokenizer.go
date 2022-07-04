package gse

import (
	"errors"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/go-ego/gse"
)

type GseTokenizer struct {
	seg *gse.Segmenter
}

// NewGseTokenizer create a gse cut tokenizer
func NewGseTokenizer(dictPath, stopPath string, alpha bool) (*GseTokenizer, error) {
	var (
		seg gse.Segmenter
		err error
	)

	seg.SkipLog = true
	if alpha {
		seg.AlphaNum = true
	}

	if dictPath != "" {
		err = seg.LoadDict(dictPath)
		if err != nil {
			return nil, err
		}
	} else {
		err = seg.LoadDict()
		if err != nil {
			return nil, err
		}
	}

	if stopPath != "" {
		err = seg.LoadStop(stopPath)
		if err != nil {
			return nil, err
		}
	}
	return &GseTokenizer{seg: &seg}, nil
}

func (g *GseTokenizer) Free() {}

// Tokenize cut the text to bleve token stream
func (g *GseTokenizer) Tokenize(text []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	t1 := string(text)
	cuts := g.seg.Trim(g.seg.CutSearch(t1, true))
	// fmt.Println("cuts: ", cuts)
	azs := g.seg.Analyze(cuts, t1)
	for _, az := range azs {
		token := analysis.Token{
			Term:     []byte(az.Text),
			Start:    az.Start,
			End:      az.End,
			Position: az.Position,
			Type:     analysis.Ideographic,
		}
		result = append(result, &token)
	}
	return result
}

func tokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	dictpath, ok := config["dictpath"].(string)
	if !ok {
		return nil, errors.New("config dictpath not found")
	}
	stoppath, ok := config["stoppath"].(string)
	if !ok {
		return nil, errors.New("config stoppath not found")
	}
	alpha, ok := config["alpha"].(bool)
	if !ok {
		return nil, errors.New("config alpha not found")
	}
	return NewGseTokenizer(dictpath, stoppath, alpha)
}

func init() {
	registry.RegisterTokenizer("gse", tokenizerConstructor)
}

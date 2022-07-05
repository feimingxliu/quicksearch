package gse

import (
	"errors"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/go-ego/gse"
)

const Name = "gse"

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
		err = seg.LoadDictEmbed()
		if err != nil {
			return nil, err
		}
	}

	if stopPath != "" {
		err = seg.LoadStop(stopPath)
		if err != nil {
			return nil, err
		}
	} else {
		err = seg.LoadStopEmbed()
		if err != nil {
			return nil, err
		}
	}
	return &GseTokenizer{seg: &seg}, nil
}

// Tokenize cut the text to bleve token stream
func (g *GseTokenizer) Tokenize(text []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	t := string(text)
	cuts := g.seg.Trim(g.seg.CutSearch(t, true))
	// fmt.Println("cuts: ", cuts)
	azs := g.seg.Analyze(cuts, t)
	for _, az := range azs {
		typ := analysis.Ideographic
		alphaNumeric := true
		for _, r := range az.Text {
			if r < 32 || r > 126 {
				alphaNumeric = false
				break
			}
		}
		if alphaNumeric {
			typ = analysis.AlphaNumeric
		}
		token := analysis.Token{
			Term:     []byte(az.Text),
			Start:    az.Start,
			End:      az.End,
			Position: az.Position,
			Type:     typ,
		}
		result = append(result, &token)
	}
	return result
}

func tokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	dictpath, ok := config["dict_path"].(string)
	if !ok {
		return nil, errors.New("config dict_path not found")
	}
	stoppath, ok := config["stop_words"].(string)
	if !ok {
		return nil, errors.New("config stop_words not found")
	}
	alpha, ok := config["alpha"].(bool)
	if !ok {
		return nil, errors.New("config alpha not found")
	}
	return NewGseTokenizer(dictpath, stoppath, alpha)
}

func init() {
	registry.RegisterTokenizer(Name, tokenizerConstructor)
}

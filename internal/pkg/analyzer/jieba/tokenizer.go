package jieba

import (
	"errors"
	"path"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/yanyiwu/gojieba"
)

const Name = "gojieba"

type JiebaTokenizer struct {
	handle *gojieba.Jieba
}

func NewJiebaTokenizer(dictpath, hmmpath, userdictpath, idf, stop_words string) *JiebaTokenizer {
	if dictpath == "" {
		dictpath = path.Join(DictDir, "jieba.dict.utf8")
	}
	if hmmpath == "" {
		hmmpath = path.Join(DictDir, "hmm_model.utf8")
	}
	if userdictpath == "" {
		userdictpath = path.Join(DictDir, "user.dict.utf8")
	}
	if idf == "" {
		idf = path.Join(DictDir, "idf.utf8")
	}
	if stop_words == "" {
		stop_words = path.Join(DictDir, "stop_words.utf8")
	}
	x := gojieba.NewJieba(dictpath, hmmpath, userdictpath, idf, stop_words)
	return &JiebaTokenizer{x}
}

func (x *JiebaTokenizer) Free() {
	x.handle.Free()
}

func (x *JiebaTokenizer) Tokenize(sentence []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	pos := 1
	words := x.handle.Tokenize(string(sentence), gojieba.SearchMode, true)
	for _, word := range words {
		typ := analysis.Ideographic
		alphaNumeric := true
		for _, r := range word.Str {
			if r < 32 || r > 126 {
				alphaNumeric = false
				break
			}
		}
		if alphaNumeric {
			typ = analysis.AlphaNumeric
		}
		token := analysis.Token{
			Term:     []byte(word.Str),
			Start:    word.Start,
			End:      word.End,
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
	hmmpath, ok := config["hmm_path"].(string)
	if !ok {
		return nil, errors.New("config hmm_path not found")
	}
	userdictpath, ok := config["user_dict_path"].(string)
	if !ok {
		return nil, errors.New("config user_dict_path not found")
	}
	idf, ok := config["idf"].(string)
	if !ok {
		return nil, errors.New("config idf not found")
	}
	stop_words, ok := config["stop_words"].(string)
	if !ok {
		return nil, errors.New("config stop_words not found")
	}
	return NewJiebaTokenizer(dictpath, hmmpath, userdictpath, idf, stop_words), nil
}

func init() {
	registry.RegisterTokenizer(Name, tokenizerConstructor)
}

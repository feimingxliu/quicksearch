package jieba

import (
	"github.com/feimingxliu/quicksearch/pkg/tokenizer"
	"github.com/yanyiwu/gojieba"
)

func NewJieBa() *jieba {
	return &jieba{jb: gojieba.NewJieba()}
}

type jieba struct {
	jb *gojieba.Jieba
}

//Tokenize return all tokens of s.
func (j *jieba) Tokenize(s string) []string {
	return j.jb.CutForSearch(s, true)
}

//Keywords return s's keywords.
func (j *jieba) Keywords(s string, topK int) []string {
	return j.jb.Extract(s, topK)
}

//KeywordsWeight return s's keywords with weight.
func (j *jieba) KeywordsWeight(s string, topK int) []tokenizer.WordWeight {
	ww := j.jb.ExtractWithWeight(s, topK)
	w := make([]tokenizer.WordWeight, len(ww))
	for i := range ww {
		w[i] = tokenizer.WordWeight{
			Word:   ww[i].Word,
			Weight: ww[i].Weight,
		}
	}
	return w
}

//Close release the resources.
func (j *jieba) Close() {
	j.jb.Free()
}

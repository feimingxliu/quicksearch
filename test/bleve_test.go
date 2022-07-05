package test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	bindex "github.com/blevesearch/bleve_index_api"
	_ "github.com/feimingxliu/quicksearch/internal/pkg/analyzer/gse"
	"github.com/feimingxliu/quicksearch/internal/pkg/analyzer/jieba"
	_ "github.com/feimingxliu/quicksearch/internal/pkg/analyzer/sego"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/yanyiwu/gojieba"
	"log"
	"os"
	"testing"
)

const (
	docsFile  = "testdata/zhwiki-20220601-abstract.json"
	numOfDoc  = 1000
	indexPath = "data/zhwiki"
)

type document struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Text  string `json:"text"`
}

func indexSomeDocs(t *testing.T, index bleve.Index) {
	f, err := os.OpenFile(docsFile, os.O_RDONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	batch := index.NewBatch()
	batchSize := 1000
	totalBulked := 0
	duration := util.ExecTime(func() {
		currentSize := 0
		for i := 0; i < numOfDoc && scanner.Scan(); i++ {
			doc := new(document)
			if err := json.Unmarshal(scanner.Bytes(), doc); err != nil {
				t.Fatal(err)
			}
			err := batch.Index(doc.ID, doc)
			if err != nil {
				t.Fatal(err)
			}
			totalBulked++
			currentSize++
			if currentSize >= batchSize {
				err = index.Batch(batch)
				if err != nil {
					t.Fatal(err)
				}
				batch.Reset()
				currentSize = 0
			}
		}
		// batch remaining
		err := index.Batch(batch)
		if err != nil {
			t.Fatal(err)
		}
	})
	log.Printf("Index %d docs costs: %s\n", totalBulked, duration)
}

// go test -v github.com/feimingxliu/quicksearch/test -run 'Bleve$'
func TestBleve(t *testing.T) {
	err := os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	// construct mapping
	mapping := bleve.NewIndexMapping()
	wikiMapping := bleve.NewDocumentMapping()
	idFieldMapping := bleve.NewKeywordFieldMapping()
	//idFieldMapping.Store = false  // not work
	wikiMapping.AddFieldMappingsAt("id", idFieldMapping)
	titleFieldMapping := bleve.NewTextFieldMapping()
	//titleFieldMapping.Store = false
	wikiMapping.AddFieldMappingsAt("title", titleFieldMapping)
	urlFieldMapping := bleve.NewTextFieldMapping()
	//urlFieldMapping.Store = false
	wikiMapping.AddFieldMappingsAt("url", urlFieldMapping)
	textFieldMapping := bleve.NewTextFieldMapping()
	wikiMapping.AddFieldMappingsAt("text", textFieldMapping)

	// open the index
	index, err := bleve.New(indexPath, mapping)
	//index, err := bleve.NewUsing(indexPath, mapping, upsidedown.Name, bleve.Config.DefaultKVStore, nil)
	if err != nil {
		t.Fatal(err)
	}

	// bulk some docs
	indexSomeDocs(t, index)

	// search
	query := bleve.NewMatchQuery(`数学，是研究数量、结构以及空间等概念及其变化的一门学科，从某种角度看属於形式科学的一种`)
	query.SetBoost(1)
	query.FieldVal = "text"
	search := bleve.NewSearchRequest(query)
	// this will make search return the doc source in "fields"
	search.Fields = []string{"*"}
	search.Size = 1
	// highlight the result
	//search.Highlight = bleve.NewHighlight()
	//search.Highlight.AddField("text")
	search.IncludeLocations = false
	searchResults, err := index.Search(search)
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	// set this to disable escape the `<` and `>`
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(searchResults)
	fmt.Printf("result:\n%s\n", buf.String())
	if len(searchResults.Hits) > 0 {
		docID := searchResults.Hits[0].ID
		doc, err := index.Document(docID)
		if err != nil {
			t.Fatal(err)
		}
		doc.VisitFields(func(field bindex.Field) {
			fmt.Printf("key: %s, value: %s\n", field.Name(), field.Value())
		})
	}

	json.Print("stats", index.StatsMap())

	err = index.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBleveWithJieba(t *testing.T) {
	err := os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	// construct mapping
	mapping := bleve.NewIndexMapping()
	err = mapping.AddCustomTokenizer("gojieba",
		map[string]interface{}{
			"dict_path":      gojieba.DICT_PATH,
			"hmm_path":       gojieba.HMM_PATH,
			"user_dict_path": gojieba.USER_DICT_PATH,
			"idf":            gojieba.IDF_PATH,
			"stop_words":     gojieba.STOP_WORDS_PATH,
			"type":           "gojieba",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	err = mapping.AddCustomAnalyzer("gojieba",
		map[string]interface{}{
			"type":      "gojieba",
			"tokenizer": "gojieba",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	mapping.DefaultAnalyzer = "gojieba"

	// open the index
	index, err := bleve.New(indexPath, mapping)
	if err != nil {
		t.Fatal(err)
	}

	// bulk some docs
	indexSomeDocs(t, index)

	// search
	query := bleve.NewMatchQuery(`数学，是研究数量、结构以及空间等概念及其变化的一门学科，从某种角度看属於形式科学的一种`)
	query.SetBoost(1)
	query.FieldVal = "text"
	search := bleve.NewSearchRequest(query)
	// this will make search return the doc source in "fields"
	search.Fields = []string{"*"}
	search.Size = 1
	// highlight the result
	//search.Highlight = bleve.NewHighlight()
	//search.Highlight.AddField("text")
	search.IncludeLocations = false
	searchResults, err := index.Search(search)
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	// set this to disable escape the `<` and `>`
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(searchResults)
	fmt.Printf("result:\n%s\n", buf.String())

	// cleanup cgo allocated heap memory
	if j, ok := (index.Mapping().AnalyzerNamed("gojieba").Tokenizer).(*jieba.JiebaTokenizer); !ok {
		t.Fatal("jieba.Free() failed")
	} else {
		j.Free()
	}

	err = index.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBleveWithGse(t *testing.T) {
	err := os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	// construct mapping
	mapping := bleve.NewIndexMapping()
	err = mapping.AddCustomTokenizer("gse",
		map[string]interface{}{
			"type":       "gse",
			"dict_path":  "",
			"stop_words": "",
			"alpha":      false,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	err = mapping.AddCustomAnalyzer("gse",
		map[string]interface{}{
			"type":      "gse",
			"tokenizer": "gse",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	mapping.DefaultAnalyzer = "gse"

	// open the index
	index, err := bleve.New(indexPath, mapping)
	if err != nil {
		t.Fatal(err)
	}

	// bulk some docs
	indexSomeDocs(t, index)

	// search
	query := bleve.NewMatchQuery(`数学，是研究数量、结构以及空间等概念及其变化的一门学科，从某种角度看属於形式科学的一种`)
	query.SetBoost(1)
	query.FieldVal = "text"
	search := bleve.NewSearchRequest(query)
	// this will make search return the doc source in "fields"
	search.Fields = []string{"*"}
	search.Size = 1
	// highlight the result
	//search.Highlight = bleve.NewHighlight()
	//search.Highlight.AddField("text")
	search.IncludeLocations = false
	searchResults, err := index.Search(search)
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	// set this to disable escape the `<` and `>`
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(searchResults)
	fmt.Printf("result:\n%s\n", buf.String())

	err = index.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBleveWithSego(t *testing.T) {
	err := os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	// construct mapping
	mapping := bleve.NewIndexMapping()
	err = mapping.AddCustomTokenizer("sego",
		map[string]interface{}{
			"type":      "sego",
			"dict_path": "",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	err = mapping.AddCustomAnalyzer("sego",
		map[string]interface{}{
			"type":      "sego",
			"tokenizer": "sego",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	mapping.DefaultAnalyzer = "sego"

	// open the index
	index, err := bleve.New(indexPath, mapping)
	if err != nil {
		t.Fatal(err)
	}

	// bulk some docs
	indexSomeDocs(t, index)

	// search
	query := bleve.NewMatchQuery(`数学，是研究数量、结构以及空间等概念及其变化的一门学科，从某种角度看属於形式科学的一种`)
	query.SetBoost(1)
	query.FieldVal = "text"
	search := bleve.NewSearchRequest(query)
	// this will make search return the doc source in "fields"
	search.Fields = []string{"*"}
	search.Size = 1
	// highlight the result
	//search.Highlight = bleve.NewHighlight()
	//search.Highlight.AddField("text")
	search.IncludeLocations = false
	searchResults, err := index.Search(search)
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	// set this to disable escape the `<` and `>`
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(searchResults)
	fmt.Printf("result:\n%s\n", buf.String())

	err = index.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}
}

package test

import (
	"bufio"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	bindex "github.com/blevesearch/bleve_index_api"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"log"
	"os"
	"testing"
)

const (
	docsFile  = "testdata/zhwiki-20220601-abstract.json"
	numOfDoc  = 10000
	indexPath = "testdata/zhwiki"
)

type document struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Text  string `json:"text"`
}

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
			err = batch.Index(doc.ID, doc)
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
		err = index.Batch(batch)
		if err != nil {
			t.Fatal(err)
		}
	})
	log.Printf("Index %d docs costs: %s\n", totalBulked, duration)

	// search
	query := bleve.NewMatchQuery(`数学，是研究数量、结构以及空间等概念及其变化的一门学科，从某种角度看属於形式科学的一种`)
	query.SetBoost(1)
	query.FieldVal = "text"
	search := bleve.NewSearchRequest(query)
	// this will make search return the doc source int "fields"
	search.Fields = []string{"*"}
	//search.Size = 1
	//search.Highlight = bleve.NewHighlight()
	//search.Highlight.AddField("text")
	search.IncludeLocations = false
	searchResults, err := index.Search(search)
	if err != nil {
		t.Fatal(err)
	}

	json.Print("result", searchResults)
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

	err = index.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}
}

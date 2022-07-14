package core

import (
	"bufio"
	"github.com/blevesearch/bleve/v2"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"log"
	"os"
	"testing"
	"time"
)

const (
	docsFile   = "../../test/testdata/zhwiki-20220601-abstract.json"
	docMapping = `{
    "default_mapping": {
        "properties": {
            "id": {
                "fields": [
                    {
                        "type": "keyword"
                    }
                ]
            },
            "title": {
                "fields": [
                    {
                        "type": "text"
                    }
                ]
            },
            "url": {
                "disabled": true,
                "fields": [
                    {
                        "type": "keyword"
                    }
                ]
            },
            "text": {
                "fields": [
                    {
                        "type": "text"
                    }
                ]
            }
        }
    },
    "type_field": "_type",
    "default_type": "_default",
    "default_analyzer": "default"
}`
)

func indexSomeDocs(t *testing.T, num int) {
	// build index mapping
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(docMapping), &m)
	if err != nil {
		t.Fatal(err)
	}
	im, err := BuildIndexMappingFromMap(m)
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewIndex(WithName(indexName), WithIndexMapping(im), WithShards(1))
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(docsFile, os.O_RDONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	totalBulked := 0
	duration := util.ExecTime(func() {
		for i := 0; i < num && scanner.Scan(); i++ {
			doc := make(map[string]interface{})
			if err := json.Unmarshal(scanner.Bytes(), &doc); err != nil {
				t.Fatal(err)
			}
			if i == 0 {
				firstDocID = doc["id"].(string)
			}
			err := index.IndexOrUpdateDocument(doc["id"].(string), doc)
			if err != nil {
				t.Fatal(err)
			}
			totalBulked++
		}
	})
	log.Printf("Index %d docs costs: %s\n", totalBulked, duration)
	err = index.Close()
	if err != nil {
		t.Fatal(err)
	}
}

var firstDocID = `43a1b5cd7383441b83049dc85188d9f3`

func TestIndexDocument(t *testing.T) {
	prepare(t)
	defer clean(t)
	indexSomeDocs(t, 100)

	index, err := GetIndex(indexName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		log.Println("Delete Index.")
		if err := index.Delete(); err != nil {
			t.Fatal(err)
		}
	}()

	// cause there is only one shard for this index
	bleveIndex := index.Shards[0].Indexer

	// search keyword

	query := bleve.NewTermQuery(firstDocID)
	query.SetBoost(1)
	//query.FieldVal = "id"
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.Size = 1
	search.IncludeLocations = false
	searchResults, err := bleveIndex.Search(search)
	if err != nil {
		t.Fatal(err)
	}
	json.Print("search term id", searchResults)
	// by default, the results are sorted by `score`
	if len(searchResults.Hits) > 0 {
		if searchResults.Hits[0].ID != firstDocID {
			t.Errorf("search term id [%s] failed", firstDocID)
		}
	} else {
		t.Errorf("search term id [%s] failed", firstDocID)
	}

	// search disable field
	url := "数学"
	mquery := bleve.NewMatchQuery(url)
	mquery.SetBoost(1)
	mquery.SetField("url")
	search = bleve.NewSearchRequest(mquery)
	search.Fields = []string{"*"}
	search.Size = 1
	search.IncludeLocations = false
	searchResults, err = bleveIndex.Search(search)
	if err != nil {
		t.Fatal(err)
	}
	json.Print("search term url(disabled)", searchResults)
	if len(searchResults.Hits) != 0 {
		t.Errorf("search term url [%s] failed", url)
	}

	// search match text
	text := "研究数量、结构以及空间等概念及其变化的一门学科"
	mquery = bleve.NewMatchQuery(text)
	mquery.SetBoost(1)
	//mquery.FieldVal = "text"
	search = bleve.NewSearchRequest(mquery)
	search.Fields = []string{"*"}
	search.Size = 1
	search.IncludeLocations = false
	searchResults, err = bleveIndex.Search(search)
	if err != nil {
		t.Fatal(err)
	}
	json.Print("search match text", searchResults)
	if len(searchResults.Hits) > 0 {
		if searchResults.Hits[0].ID != firstDocID {
			t.Errorf("search match text [%s] failed", text)
		}
	} else {
		t.Errorf("search match text [%s] failed", text)
	}

	trquery := bleve.NewDateRangeQuery(time.Now().Add(-1*time.Minute), time.Now())
	trquery.SetField("@timestamp")
	search = bleve.NewSearchRequest(trquery)
	search.Fields = []string{"*"}
	search.Size = 10
	search.IncludeLocations = false
	searchResults, err = bleveIndex.Search(search)
	if err != nil {
		t.Fatal(err)
	}
	json.Print("search date range", searchResults)
}

func TestGetDocument(t *testing.T) {
	prepare(t)
	defer clean(t)
	indexSomeDocs(t, 1)

	var err error
	index, err := GetIndex(indexName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		log.Println("Delete Index.")
		if err := index.Delete(); err != nil {
			t.Fatal(err)
		}
	}()

	var doc *Document
	duration := util.ExecTime(func() {
		doc, err = index.GetDocument(firstDocID)
		if err != nil {
			t.Errorf("index.GetDocument: %s", err)
		}
	})
	log.Printf("index.GetDocument costs: %d ns\n", duration)
	json.Print("document", doc)
}

func TestDeleteDocument(t *testing.T) {
	prepare(t)
	defer clean(t)
	indexSomeDocs(t, 1)

	var err error
	index, err := GetIndex(indexName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		log.Println("Delete Index.")
		if err := index.Delete(); err != nil {
			t.Fatal(err)
		}
	}()

	var doc *Document
	doc, err = index.GetDocument(firstDocID)
	if err != nil {
		t.Errorf("index.GetDocument: %s", err)
	}
	json.Print("before delete, document", doc)

	// delete document
	err = index.DeleteDocument(firstDocID)
	if err != nil {
		t.Errorf("index.DeleteDocument: %s", err)
	}

	doc, err = index.GetDocument(firstDocID)
	if err != errors.ErrDocumentNotFound {
		t.Errorf("index.GetDocument: %s", err)
	}
	if doc != nil && doc.Found {
		t.Error("document not deleted")
	}
	json.Print("after delete, document", doc)
}

func TestUpdateDocument(t *testing.T) {
	prepare(t)
	defer clean(t)
	indexSomeDocs(t, 1)

	var err error
	index, err := GetIndex(indexName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		log.Println("Delete Index.")
		if err := index.Delete(); err != nil {
			t.Fatal(err)
		}
	}()

	var doc *Document
	doc, err = index.GetDocument(firstDocID)
	if err != nil {
		t.Errorf("index.GetDocument: %s", err)
	}
	json.Print("before update, document", doc)

	// delete document
	err = index.UpdateDocumentPartially(firstDocID, map[string]interface{}{"text": "updated partially"})
	if err != nil {
		t.Errorf("index.UpdateDocumentPartially: %s", err)
	}

	doc, err = index.GetDocument(firstDocID)
	if err != nil {
		t.Errorf("index.GetDocument: %s", err)
	}
	if source, ok := doc.Source.(map[string]interface{}); ok {
		if source["text"] != "updated partially" {
			t.Error("document not updated")
		}
	}
	json.Print("after update, document", doc)
}

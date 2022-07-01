package core

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
	"log"
	"sync"
	"testing"
)

func TestIndexDocument(t *testing.T) {
	prepare(t)
	index := NewIndex(WithName(indexName), WithStorage("bolt"), WithTokenizer("jieba"))
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(raw), &m)
	doc := NewDocument(m)
	duration := util.ExecTime(func() {
		if err := index.IndexDocument(doc); err != nil {
			t.Fatal(err)
		} else {
			json.Print("index", index)
			json.Print("doc", doc)
		}
	})
	log.Println("IndexDocument costs: ", duration)
	query := m["text"].(string)
	for _, keyword := range index.tokenizer.Keywords(query, len(query)) {
		if ids, err := index.GetIDsByKeyword(keyword); err != nil {
			t.Error(err)
		} else {
			if !slices.ContainsStr(ids, doc.ID) {
				t.Errorf("Inverted index: keyword %q, %q not included\n", keyword, doc.ID)
			}
		}
	}
	log.Println("Delete Index.")
	if err := index.Delete(); err != nil {
		t.Fatal(err)
	}
}

func TestRetrieveDocument(t *testing.T) {
	prepare(t)
	index := NewIndex(WithName(indexName), WithStorage("bolt"), WithTokenizer("jieba"))
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(raw), &m)
	doc := NewDocument(m)
	if err := index.IndexDocument(doc); err != nil {
		t.Fatal(err)
	} else {
		json.Print("index", index)
		json.Print("doc", doc)
	}
	if doc, err := index.RetrieveDocument(doc.ID); err != nil {
		t.Fatal(err)
	} else {
		json.Print("RetrieveDocument", doc)
	}
	log.Println("Delete Index.")
	if err := index.Delete(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDocument(t *testing.T) {
	prepare(t)
	index := NewIndex(WithName(indexName), WithStorage("bolt"), WithTokenizer("jieba"))
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(raw), &m)
	doc := NewDocument(m)
	if err := index.IndexDocument(doc); err != nil {
		t.Fatal(err)
	} else {
		json.Print("index", index)
		json.Print("doc", doc)
	}
	if ids, err := index.GetIDsByKeyword("数学"); err != nil {
		t.Fatal(err)
	} else {
		json.Print("数学", ids)
	}
	if err := index.DeleteDocument(doc.ID); err != nil {
		t.Fatal(err)
	}
	if ids, err := index.GetIDsByKeyword("数学"); err != nil {
		t.Fatal(err)
	} else {
		json.Print("After DeleteDocument, search `数学`", ids)
	}
	log.Println("Delete Index.")
	if err := index.Delete(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateDocument(t *testing.T) {
	prepare(t)
	index := NewIndex(WithName(indexName), WithStorage("bolt"), WithTokenizer("jieba"))
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(raw), &m)
	doc := NewDocument(m)
	duration := util.ExecTime(func() {
		if err := index.IndexDocument(doc); err != nil {
			t.Fatal(err)
		} else {
			json.Print("index", index)
			json.Print("doc", doc)
		}
	})
	log.Println("IndexDocument costs: ", duration)
	if ids, err := index.GetIDsByKeyword("数学"); err != nil {
		t.Fatal(err)
	} else {
		json.Print("Before update, search `数学`", ids)
	}
	//change the doc source.
	doc.Source = map[string]interface{}{}
	duration = util.ExecTime(func() {
		if err := index.UpdateDocument(doc); err != nil {
			t.Fatal(err)
		} else {
			json.Print("index", index)
			json.Print("doc", doc)
		}
	})
	log.Println("UpdateDocument costs: ", duration)
	if ids, err := index.GetIDsByKeyword("数学"); err != nil {
		t.Fatal(err)
	} else {
		json.Print("After update, search `数学`", ids)
	}
	log.Println("Delete Index.")
	if err := index.Delete(); err != nil {
		t.Fatal(err)
	}
}

func TestBulkIndexDocument1000(t *testing.T) {
	bulkIndexDocument(t, 1000)
}

//go test -v -timeout 0 github.com/feimingxliu/quicksearch/internal/core -run 'BulkIndexDocument10000'
func TestBulkIndexDocument10000(t *testing.T) {
	bulkIndexDocument(t, 10000)
}

func bulkIndexDocument(t *testing.T, n uint) {
	prepare(t)
	index := NewIndex(WithName(indexName), WithStorage("bolt"), WithTokenizer("jieba"))
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(raw), &m)
	docs := make([]*Document, n)
	for i := range docs {
		docs[i] = NewDocument(m)
	}
	duration := util.ExecTime(func() {
		var wg sync.WaitGroup
		batchSize := 1000               //one batch size
		pieces := len(docs) / batchSize //divided into pieces
		for k := 0; k < pieces; k++ {
			wg.Add(1)
			go func(k int) {
				if err := index.BulkDocuments(docs[k*batchSize : (k+1)*batchSize]); err != nil {
					fmt.Println(err)
				}
				log.Printf("piece %d done!\n", k)
				wg.Done()
			}(k)
		}
		//bulk remaining
		if len(docs)%batchSize > 0 {
			if err := index.BulkDocuments(docs[pieces*batchSize:]); err != nil {
				fmt.Println(err)
			}
			log.Println("remaining piece done!")
		}
		wg.Wait()
	})
	log.Printf("Bulk %d documents, costs: %s\n", len(docs), duration)
	if err := index.UpdateMetadata(); err != nil {
		t.Fatal(err)
	}
	//check if inverted index works.
	if ids, err := index.GetIDsByKeyword("数学"); err != nil {
		t.Fatal(err)
	} else {
		json.Print("数学", ids)
	}
	log.Println("Delete Index.")
	if err := index.Delete(); err != nil {
		t.Fatal(err)
	}
}

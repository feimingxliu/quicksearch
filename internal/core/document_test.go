package core

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
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
	if err := index.IndexDocument(doc); err != nil {
		t.Fatal(err)
	} else {
		json.Print("index", index)
		json.Print("doc", doc)
	}
	if ids, err := index.GetIDsByKeyword("数学"); err != nil {
		t.Fatal(err)
	} else {
		json.Print("Before update, search `数学`", ids)
	}
	//change the doc source.
	doc.Source = map[string]interface{}{}
	if err := index.IndexDocument(doc); err != nil {
		t.Fatal(err)
	} else {
		json.Print("index same doc", index)
		json.Print("doc", doc)
	}
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

/*func TestBulkIndexDocument(t *testing.T) {
	prepare(t)
	index := NewIndex(WithName(indexName), WithStorage("bolt"), WithTokenizer("jieba"))
	f, err := os.OpenFile("../../test/testdata/zhwiki-20220601-abstract.json", os.O_RDONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(f)
	docsRaw := make([]map[string]interface{}, 0)
	log.Println("Decoding......")
	duration := util.ExecTime(func() {
		if err := decoder.Decode(&docsRaw); err != nil {
			t.Fatal(err)
		}
	})
	log.Println("Decoding json file costs: ", duration)
	docs := make([]*Document, 10000, 10000)
	for i := 0; i < len(docs); i++ {
		docs[i] = NewDocument(docsRaw[i])
		docs[i].WithID(docsRaw[i]["id"].(string))
	}
	duration = util.ExecTime(func() {
		var wg sync.WaitGroup
		pieces := 100 //divided into pieces
		base := len(docs) / pieces
		for k := 0; k < pieces; k++ {
			wg.Add(1)
			go func(k int) {
				if err := index.BulkDocuments(docs[k*base : (k+1)*base]); err != nil {
					fmt.Println(err)
				}
				log.Printf("piece %d done!", k)
				wg.Done()
			}(k)
		}
		wg.Wait()
	})
	log.Printf("Bulk %d documents, costs: %s\n", len(docs), duration)
	if err := index.UpdateMetadata(); err != nil {
		t.Fatal(err)
	}
	//if ids, err := index.GetIDsByKeyword("数学"); err != nil {
	//	t.Fatal(err)
	//} else {
	//	json.Print("数学", ids)
	//}
	log.Println("Close Index.")
	if err := index.Close(); err != nil {
		t.Fatal(err)
	}
}*/

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

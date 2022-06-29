package core

import (
	"bufio"
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

//go test -v -timeout 0 github.com/feimingxliu/quicksearch/internal/core  -run 'Search$' -memprofile mem.out
//pprof -http=:8080 mem.out
func TestSearch(t *testing.T) {
	prepare(t)
	index := NewIndex(WithName(indexName), WithStorage("bolt"), WithTokenizer("jieba"))
	if index.DocNum == 0 {
		indexSomeDocs(t)
		// wait batches to trigger
		time.Sleep(30 * time.Second)
	}
	json.Print("index", index)
	var res *SearchResult
	duration := util.ExecTime(func() {
		res = index.Search(
			NewSearchOption().
				SetQuery("数学，是研究数量、结构以及空间等概念及其变化的一门学科，从某种角度看属於形式科学的一种。").
				SetTopN(10).SetTimeout(1 * time.Second),
		)
	})
	log.Println("Search costs: ", duration)
	json.Print("SearchResult", res)
	log.Println("Delete Index.")
	if err := index.Delete(); err != nil {
		t.Fatal(err)
	}
}

func indexSomeDocs(t *testing.T) {
	index := NewIndex(WithName(indexName), WithStorage("bolt"), WithTokenizer("jieba"))
	f, err := os.OpenFile("../../test/testdata/zhwiki-20220601-abstract.json", os.O_RDONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	const numOfDoc = 10000
	docs := make([]*Document, 0, numOfDoc)
	scanner := bufio.NewScanner(f)
	duration := util.ExecTime(func() {
		for i := 0; i < numOfDoc && scanner.Scan(); i++ {
			m := make(map[string]interface{})
			if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
				t.Fatal(err)
			}
			docs = append(docs, NewDocument(m).WithID(m["id"].(string)))
		}
	})
	log.Println("Parsing json file costs: ", duration)
	duration = util.ExecTime(func() {
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
}

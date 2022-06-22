package core

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"log"
	"sync"
	"testing"
)

const indexName = "test"

var once sync.Once

func prepare(t *testing.T) {
	//use once, so it will only init once.
	once.Do(func() {
		if err := config.Init("../../configs/config.yaml"); err != nil {
			t.Fatal(err)
		}
		//change DataDir, because `pwd` is in the source code dir.
		config.Global.Storage.DataDir = "../../data"
		InitMeta()
	})
}

func TestIndex(t *testing.T) {
	prepare(t)
	index := NewIndex(WithName(indexName), WithStorage("bolt"), WithTokenizer("jieba"))
	if err := index.UpdateMetadata(); err != nil {
		t.Fatal(err)
	}
	if index, err := GetIndex(indexName); err != nil {
		t.Fatal(err)
	} else {
		json.Print("GetIndex", index)
	}
	if indexes, err := ListIndices(); err != nil {
		t.Fatal(err)
	} else {
		json.Print("ListIndexes", indexes)
	}
	if err := index.Close(); err != nil {
		t.Fatal(err)
	}
	if err := index.Open(); err != nil {
		t.Fatal(err)
	}
	if err := index.Delete(); err != nil {
		t.Fatal(err)
	} else {
		if _, err := GetIndex(indexName); err != errors.ErrKeyNotFound {
			t.Fatal("DeleteIndex failed, index not deleted!")
		}
	}
}

func TestCloneIndex(t *testing.T) {
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
	cloneName := indexName + "-clone"
	if err := index.Clone(cloneName); err != nil {
		t.Fatal(err)
	}
	clone, err := GetIndex(cloneName)
	if err != nil {
		t.Fatal(err)
	} else {
		json.Print("clone index", clone)
	}
	if ids, err := clone.GetIDsByKeyword("数学"); err != nil {
		t.Fatal(err)
	} else {
		json.Print("数学", ids)
	}
	log.Println("Delete Index.")
	if err := index.Delete(); err != nil {
		t.Fatal(err)
	}
	if err := clone.Delete(); err != nil {
		t.Fatal(err)
	}
}

const raw = `{
  "id": "cc5caccfdcc743c387cfacce7ebb0369",
  "title": "Wikipedia: 数学",
  "url": "https://zh.wikipedia.org/wiki/数学",
  "text": "数学，是研究数量、结构以及空间等概念及其变化的一门学科，从某种角度看属于形式科学的一种。"
 }`

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

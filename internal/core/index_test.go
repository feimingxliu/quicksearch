package core

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"sync"
	"testing"
	"time"
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
		Init()
	})
}

func TestIndex(t *testing.T) {
	prepare(t)
	index := NewIndex(indexName)
	if err := index.UpdateMetadata(); err != nil {
		t.Fatal(err)
	}
	if index, err := GetIndex(indexName); err != nil {
		t.Fatal(err)
	} else {
		json.Print("GetIndex", index)
	}
	if indexes, err := ListIndexes(); err != nil {
		t.Fatal(err)
	} else {
		json.Print("ListIndexes", indexes)
	}
	if err := DeleteIndex(indexName); err != nil {
		t.Fatal(err)
	} else {
		if index, err := GetIndex(indexName); err != nil {
			t.Fatal(err)
		} else {
			if index != nil {
				t.Fatal("DeleteIndex failed, index not deleted!")
			}
		}
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
	index := NewIndex(indexName)
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(raw), &m)
	doc := NewDocument(m)
	if err := index.IndexDocument(doc); err != nil {
		t.Fatal(err)
	} else {
		json.Print("index", index)
		json.Print("doc", doc)
	}
	pairs, err := InvertedIndex.listAllKeywordIDs()
	if err != nil {
		t.Fatal(err)
	}
	json.Print("inverted index", pairs)
}

func TestBulkDocument(t *testing.T) {
	prepare(t)
	index := NewIndex(indexName)
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(raw), &m)
	docs := make([]*Document, 1000, 1000)
	for i := 0; i < 1000; i++ {
		docs[i] = NewDocument(m)
	}
	var wg sync.WaitGroup
	for k := 0; k < 10; k++ {
		go func(k int) {
			wg.Add(1)
			if err := index.BulkDocuments(docs[k*100 : k*100+99]); err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(k)
	}
	wg.Wait()
	//wait underlying db to execute batch.
	time.Sleep(10 * time.Second)
	pairs, err := InvertedIndex.listAllKeywordIDs()
	if err != nil {
		t.Fatal(err)
	}
	json.Print("inverted index", pairs)
}

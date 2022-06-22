package core

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/errors"
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

const raw = `{
  "id": "cc5caccfdcc743c387cfacce7ebb0369",
  "title": "Wikipedia: 数学",
  "url": "https://zh.wikipedia.org/wiki/数学",
  "text": "数学，是研究数量、结构以及空间等概念及其变化的一门学科，从某种角度看属于形式科学的一种。"
 }`

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
		if _, err := GetIndex(indexName); err != errors.ErrIndexNotFound {
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

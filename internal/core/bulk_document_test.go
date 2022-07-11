package core

import (
	"bufio"
	"bytes"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"runtime"
	"testing"
)

func TestBulk(t *testing.T) {
	prepare(t)
	defer clean(t)
	bulkDocument(t, 10000)
}

func bulkDocument(t *testing.T, num int) {
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
	index, err := NewIndex(WithName(indexName), WithIndexMapping(im))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		log.Println("Delete Index.")
		if err := index.Delete(); err != nil {
			t.Fatal(err)
		}
	}()
	f, err := os.OpenFile(docsFile, os.O_RDONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(f)

	// build bulk
	bulk := &bytes.Buffer{}
	for i := 0; scanner.Scan() && i < num; i++ {
		bulk.WriteString(`{"index": {}}`)
		bulk.WriteByte('\n')
		bulk.Write(scanner.Bytes())
		if i < num-1 {
			bulk.WriteByte('\n')
		}
	}

	res, err := Bulk(indexName, bulk)
	if err != nil {
		t.Errorf("Bulk: %s", err)
	}
	log.Printf("Bulk %d docs costs: %s\n", num, res.Took)
	err = index.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBulkIndexDocument10000(t *testing.T) {
	prepare(t)
	defer clean(t)
	bulkIndexDocument(t, 10000)
}

func TestBulkIndexDocument100000(t *testing.T) {
	prepare(t)
	defer clean(t)
	bulkIndexDocument(t, 100000)
}

//go test -v -timeout 0 github.com/feimingxliu/quicksearch/internal/core -run 'BulkIndexDocument10000'  -memprofile mem.out
func TestBulkIndexDocument1000000(t *testing.T) {
	prepare(t)
	defer clean(t)
	bulkIndexDocument(t, 1000000)
}

func bulkIndexDocument(t *testing.T, num int) {
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
	index, err := NewIndex(WithName(indexName), WithIndexMapping(im))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		log.Println("Delete Index.")
		if err := index.Delete(); err != nil {
			t.Fatal(err)
		}
	}()
	f, err := os.OpenFile(docsFile, os.O_RDONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	totalBulked := 0
	wg := &errgroup.Group{}
	// if don't set limit, the following will produce thousands of goroutines when `num` is big which may run out of memory.
	wg.SetLimit(runtime.NumCPU())
	duration := util.ExecTime(func() {
		batchSize := 1000
		docs := make([]map[string]interface{}, 0, batchSize)
		for i := 0; i < num && scanner.Scan(); i++ {
			doc := make(map[string]interface{})
			if err = json.Unmarshal(scanner.Bytes(), &doc); err != nil {
				t.Fatal(err)
			}
			docs = append(docs, doc)
			if len(docs) >= batchSize {
				cdocs := docs
				wg.Go(func() error {
					err = index.BulkIndex(cdocs)
					if err != nil {
						t.Errorf("index.BulkIndex: %s", err)
						return err
					}
					totalBulked += len(cdocs)
					return nil
				})

				docs = make([]map[string]interface{}, 0, batchSize)
			}
		}
		err = wg.Wait()
		if err != nil {
			t.Errorf("errorgroup.Wait: %s", err)
		}
		err = index.BulkIndex(docs)
		if err != nil {
			t.Errorf("index.BulkIndex: %s", err)
		}
		totalBulked += len(docs)
	})
	log.Printf("Bulk %d docs costs: %s\n", totalBulked, duration)
	err = index.Close()
	if err != nil {
		t.Fatal(err)
	}
}

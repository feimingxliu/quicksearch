package test

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"github.com/syndtr/goleveldb/leveldb"
	"path/filepath"
	"testing"
)

func TestLeveldb(t *testing.T) {
	// The returned DB instance is safe for concurrent use. Which mean that all
	// DB's methods may be called concurrently from multiple goroutine.
	db, err := leveldb.OpenFile("testdata/leveldb.db", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Get([]byte("key"), nil)
	if err != leveldb.ErrNotFound {
		t.Fatal(err)
	}
	err = db.Put([]byte("key"), []byte("value"), nil)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Delete([]byte("key"), nil)
	if err != nil {
		t.Fatal(err)
	}

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("key: %s, value: %s\n", key, value)
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		t.Fatal(err)
	}

	batch := new(leveldb.Batch)
	batch.Put([]byte("foo"), []byte("value"))
	batch.Put([]byte("bar"), []byte("another value"))
	batch.Delete([]byte("baz"))
	err = db.Write(batch, nil)
	if err != nil {
		t.Fatal(err)
	}

	// try to copy db file
	// must close db first
	db.Close()

	src, _ := filepath.Abs("testdata/leveldb.db")
	err = util.CopyDir(src, "testdata/leveldb_copy.db")
	if err != nil {
		t.Fatal(err)
	}

	db, err = leveldb.OpenFile("testdata/leveldb_copy.db", nil)
	if err != nil {
		t.Fatal(err)
	}
}

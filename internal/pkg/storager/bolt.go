package storager

import (
	"bytes"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"go.etcd.io/bbolt"
	"log"
	"os"
	"path"
)

func newBolt(dbPath string) *bolt {
	opt := &bbolt.Options{
		Timeout:      0,
		NoGrowSync:   false,
		FreelistType: bbolt.FreelistArrayType,
	}
	if err := os.MkdirAll(path.Dir(dbPath), 0755); err != nil {
		log.Fatalln("[newBolt] os.MkdirAll: ", err)
	}
	db, err := bbolt.Open(dbPath, 0666, opt)
	if err != nil {
		log.Fatalln("[newBolt] bbolt.Open: ", err)
	}
	return &bolt{db}
}

type bolt struct {
	db *bbolt.DB
}

/*
Index bucket name: /index
	So set a index: bucket.Put("index_name", value)
Doc bucket name: /index/index_name
	So set a doc: bucket.Put("docID", value)
*/

//List lists all items in a bucket, so it can list indexes and docs.
func (b *bolt) List(s string) ([][]byte, error) {
	data := make([][]byte, 0)
	bucket, _ := b.splitBucketAndKey(s)
	err := b.db.View(func(txn *bbolt.Tx) error {
		b := txn.Bucket(bucket)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			valCopy := make([]byte, len(v))
			copy(valCopy, v)
			data = append(data, valCopy)
		}
		return nil
	})
	return data, err
}

func (b *bolt) Get(key string) ([]byte, error) {
	var data []byte
	bucket, name := b.splitBucketAndKey(key)
	err := b.db.View(func(txn *bbolt.Tx) error {
		b := txn.Bucket(bucket)
		if b == nil {
			return errors.ErrKeyNotFound
		}
		v := b.Get(name)
		if v == nil {
			return errors.ErrKeyNotFound
		}
		data = make([]byte, len(v))
		copy(data, v)
		return nil
	})
	return data, err
}

func (b *bolt) Set(key string, value []byte) error {
	if key == "" {
		return errors.ErrEmptyKey
	}
	bucket, name := b.splitBucketAndKey(key)
	return b.db.Update(func(txn *bbolt.Tx) error {
		b, err := txn.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return b.Put(name, value)
	})
}

func (b *bolt) Delete(key string) error {
	if key == "" {
		return errors.ErrEmptyKey
	}
	bucket, name := b.splitBucketAndKey(key)
	return b.db.Update(func(Tx *bbolt.Tx) error {
		b := Tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		return b.Delete(name)
	})
}

//DeleteAll deletes a bucket, it should just be used to delete a index.
func (b *bolt) DeleteAll(docBucket string) error {
	if docBucket == "" {
		return nil
	}
	indexBucket, indexName := b.splitBucketAndKey(docBucket)
	return b.db.Update(func(Tx *bbolt.Tx) error {
		b := Tx.Bucket(indexBucket)
		if b == nil {
			return nil
		}
		err := b.Delete(indexName)
		if err != nil {
			return err
		}
		return Tx.DeleteBucket([]byte(docBucket))
	})
}

func (b *bolt) Close() error {
	return b.db.Close()
}

//  /index/index_name/doc_uid
func (b *bolt) splitBucketAndKey(key string) ([]byte, []byte) {
	if key == "" {
		return nil, nil
	}
	p := bytes.LastIndex([]byte(key), []byte("/"))
	return []byte(key[:p]), []byte(key[p+1:])
}

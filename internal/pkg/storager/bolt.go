package storager

import (
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
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

var defaultBucket = []byte("default")

//List lists all items in bucket.
func (b *bolt) List() ([][]byte, error) {
	data := make([][]byte, 0)
	err := b.db.View(func(txn *bbolt.Tx) error {
		b := txn.Bucket(defaultBucket)
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
	err := b.db.View(func(txn *bbolt.Tx) error {
		b := txn.Bucket(defaultBucket)
		if b == nil {
			return errors.ErrKeyNotFound
		}
		v := b.Get([]byte(key))
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
	if len(key) == 0 {
		return errors.ErrEmptyKey
	}
	return b.db.Update(func(txn *bbolt.Tx) error {
		b, err := txn.CreateBucketIfNotExists(defaultBucket)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
}

func (b *bolt) Batch(keys []string, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.ErrKeyValueNotMatch
	}
	return b.db.Batch(func(tx *bbolt.Tx) error {
		for i, key := range keys {
			if len(key) == 0 {
				json.Print("keys", keys)
				return errors.ErrEmptyKey
			}
			b, err := tx.CreateBucketIfNotExists(defaultBucket)
			if err != nil {
				return err
			}
			err = b.Put([]byte(key), values[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (b *bolt) Delete(key string) error {
	if key == "" {
		return errors.ErrEmptyKey
	}
	return b.db.Update(func(Tx *bbolt.Tx) error {
		b := Tx.Bucket(defaultBucket)
		if b == nil {
			return nil
		}
		return b.Delete([]byte(key))
	})
}

//DeleteAll deletes a bucket.
func (b *bolt) DeleteAll() error {
	return b.db.Update(func(Tx *bbolt.Tx) error {
		bb := Tx.Bucket([]byte(defaultBucket))
		if bb == nil {
			return nil
		}
		return Tx.DeleteBucket(defaultBucket)
	})
}

func (b *bolt) Type() string {
	return "bolt"
}

func (b *bolt) Close() error {
	return b.db.Close()
}

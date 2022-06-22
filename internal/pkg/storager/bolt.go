package storager

import (
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"go.etcd.io/bbolt"
	"log"
	"os"
	"path"
	"runtime/debug"
)

func newBolt(dbPath string) *bolt {
	opt := &bbolt.Options{
		Timeout:      0,
		NoGrowSync:   false,
		FreelistType: bbolt.FreelistArrayType,
	}
	if err := os.MkdirAll(path.Dir(dbPath), 0755); err != nil {
		debug.PrintStack()
		log.Fatalln("[newBolt] os.MkdirAll: ", err)
	}
	db, err := bbolt.Open(dbPath, 0666, opt)
	if err != nil {
		debug.PrintStack()
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
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
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
	err := b.db.Update(func(txn *bbolt.Tx) error {
		b, err := txn.CreateBucketIfNotExists(defaultBucket)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *bolt) Batch(keys []string, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.ErrKeyValueNotMatch
	}
	return b.db.Batch(func(tx *bbolt.Tx) error {
		for i, key := range keys {
			if len(key) == 0 {
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
	err := b.db.Update(func(Tx *bbolt.Tx) error {
		b := Tx.Bucket(defaultBucket)
		if b == nil {
			return nil
		}
		return b.Delete([]byte(key))
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

//DeleteAll deletes a bucket.
func (b *bolt) DeleteAll() error {
	err := b.db.Update(func(Tx *bbolt.Tx) error {
		bb := Tx.Bucket(defaultBucket)
		if bb == nil {
			return nil
		}
		return Tx.DeleteBucket(defaultBucket)
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *bolt) CloneDatabase(path string) error {
	err := b.db.View(func(tx *bbolt.Tx) error {
		return tx.CopyFile(path, 0600)
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *bolt) Type() string {
	return "bolt"
}

func (b *bolt) Close() error {
	err := b.db.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

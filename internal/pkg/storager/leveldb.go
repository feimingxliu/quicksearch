package storager

import (
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"path"
	"runtime/debug"
)

func newLeveldb(dbPath string) *goleveldb {
	if err := os.MkdirAll(path.Dir(dbPath), 0755); err != nil {
		debug.PrintStack()
		log.Fatalln("[newLeveldb] os.MkdirAll: ", err)
	}
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		debug.PrintStack()
		log.Fatalln("[newLeveldb] leveldb.OpenFile: ", err)
	}
	return &goleveldb{db: db}
}

type goleveldb struct {
	db *leveldb.DB
}

func (l goleveldb) List() ([][]byte, error) {
	return nil, nil
}

func (l goleveldb) Get(key string) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.ErrEmptyKey
	}
	v, err := l.db.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, errors.ErrKeyNotFound
		}
		return nil, err
	}
	return v, nil
}

func (l goleveldb) Set(key string, value []byte) error {
	if len(key) == 0 {
		return errors.ErrEmptyKey
	}
	return l.db.Put([]byte(key), value, nil)
}

func (l goleveldb) Batch(keys []string, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.ErrKeyValueNotMatch
	}
	batch := new(leveldb.Batch)
	for i := range keys {
		batch.Put([]byte(keys[i]), values[i])
	}
	return l.db.Write(batch, nil)
}

func (l goleveldb) Delete(key string) error {
	if len(key) == 0 {
		return errors.ErrEmptyKey
	}
	return l.db.Delete([]byte(key), nil)
}

func (l goleveldb) DeleteAll() error {
	batch := new(leveldb.Batch)
	iter := l.db.NewIterator(nil, nil)
	for iter.Next() {
		var k []byte
		copy(k, iter.Key())
		batch.Delete(k)
	}
	iter.Release()
	return l.db.Write(batch, nil)
}

func (l goleveldb) CloneDatabase(newPath string) error {
	if err := os.MkdirAll(path.Dir(newPath), 0755); err != nil {
		return err
	}
	db, err := leveldb.OpenFile(newPath, nil)
	if err != nil {
		return err
	}
	batch := new(leveldb.Batch)
	iter := l.db.NewIterator(nil, nil)
	for iter.Next() {
		var k, v []byte
		copy(k, iter.Key())
		copy(v, iter.Value())
		batch.Put(k, v)
	}
	iter.Release()
	err = l.db.Write(batch, nil)
	if err != nil {
		return err
	}
	return db.Close()
}

func (l goleveldb) Type() string {
	return "leveldb"
}

func (l goleveldb) Close() error {
	return l.db.Close()
}

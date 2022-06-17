package storager

import "github.com/feimingxliu/quicksearch/pkg/storager/bolt"

type StorageType int

//types of storager
const (
	Bolt StorageType = iota

	Default = Bolt
)

func NewStorager(st StorageType, dbPath string) Storager {
	switch st {
	case Bolt:
		return bolt.NewBolt(dbPath)
	default:
		return bolt.NewBolt(dbPath)
	}
}

type Storager interface {
	List(string) ([][]byte, error)
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Delete(string) error
	DeleteAll(string) error
	Close() error
}

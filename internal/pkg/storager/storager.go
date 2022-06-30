package storager

type StorageType int

//TODO: support badger.
//types of storager
const (
	Bolt StorageType = iota
	Leveldb
	Default = Bolt
)

func NewStorager(st StorageType, dbPath string) Storager {
	switch st {
	case Bolt:
		return newBolt(dbPath)
	case Leveldb:
		return newLeveldb(dbPath)
	default:
		return newBolt(dbPath)
	}
}

type Storager interface {
	List() ([][]byte, error)                    // list all values
	Get(key string) ([]byte, error)             // get a value along with key
	Set(key string, value []byte) error         // set a key, value pair
	Batch(keys []string, values [][]byte) error // batch set key, value pairs
	Delete(key string) error                    // delete a key, value pair
	DeleteAll() error                           // delete all key, value pairs
	CloneDatabase(newPath string) error         // clone the database to the newPath
	Type() string                               // return the underlying type of db
	Close() error                               // close the db
}

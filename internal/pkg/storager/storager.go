package storager

type StorageType int

//TODO: support badger.
//types of storager
const (
	Bolt StorageType = iota

	Default = Bolt
)

func NewStorager(st StorageType, dbPath string) Storager {
	switch st {
	case Bolt:
		return newBolt(dbPath)
	default:
		return newBolt(dbPath)
	}
}

type Storager interface {
	List() ([][]byte, error)
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Batch([]string, [][]byte) error
	Delete(string) error
	DeleteAll() error
	CloneDatabase(string) error
	Type() string
	Close() error
}

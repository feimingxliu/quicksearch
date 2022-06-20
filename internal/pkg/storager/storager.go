package storager

type StorageType int

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
	List(string) ([][]byte, error)
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Delete(string) error
	DeleteAll(string) error
	Type() string
	Close() error
}

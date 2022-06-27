package core

import "sync"

var (
	Indices     = make(map[string]*Index) //the opened indexes
	indicesRwMu sync.RWMutex              //protect the Indices
)

func CloseIndices() error {
	for _, index := range Indices {
		if err := index.Close(); err != nil {
			return err
		}
	}
	return nil
}

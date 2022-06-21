package core

import "sync"

var (
	Indices     = make(map[string]*Index) //the opened indexes
	indicesRwMu sync.RWMutex
)

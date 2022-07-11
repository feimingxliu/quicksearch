package core

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/pkg/storager"
	"path"
	"strings"
	"sync"
)

// engine is only valid after Engine.Run()
var engine *Engine

func NewEngine() *Engine {
	return &Engine{
		indices: make(map[string]*Index),
	}
}

type Engine struct {
	indices map[string]*Index // the opened indexes
	meta    storager.Storager // metadata storage
	sync.RWMutex
}

func (e *Engine) Run() error {
	// the following call will use the global `engine`
	engine = e
	if err := e.initMeta(); err != nil {
		return err
	}
	if err := e.loadAllIndices(); err != nil {
		return err
	}
	return nil
}

func (e *Engine) Stop() error {
	e.Lock()
	if err := e.closeAllIndices(); err != nil {
		return err
	}
	if err := e.meta.Close(); err != nil {
		return err
	}
	engine = nil
	e.Unlock()
	return nil
}

func (e *Engine) addIndex(index *Index) {
	e.Lock()
	e.indices[index.Name] = index
	e.Unlock()
}

func (e *Engine) getIndex(name string) *Index {
	e.RLock()
	index := e.indices[name]
	e.RUnlock()
	return index
}

func (e *Engine) removeIndex(index *Index) {
	e.Lock()
	delete(e.indices, index.Name)
	e.Unlock()
}

func (e *Engine) initMeta() error {
	var err error
	switch strings.ToLower(config.Global.Storage.MetaType) {
	case "bolt":
		e.meta, err = storager.NewStorager(storager.Bolt, path.Join(config.Global.Storage.DataDir, "metadata", "meta"))
	default:
		e.meta, err = storager.NewStorager(storager.Default, path.Join(config.Global.Storage.DataDir, "metadata", "meta"))
	}
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) loadAllIndices() error {
	indices, err := ListIndices()
	if err != nil {
		return err
	}
	for _, index := range indices {
		if err := index.Open(); err != nil {
			return err
		}
		e.addIndex(index)
	}
	return nil
}

func (e *Engine) closeAllIndices() error {
	for _, index := range e.indices {
		if err := index.Close(); err != nil {
			return err
		}
	}
	return nil
}

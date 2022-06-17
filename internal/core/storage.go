package core

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/pkg/storager"
	"path"
)

var db storager.Storager

func init() {
	switch config.Global.Storage.Type {
	case "bolt":
		db = storager.NewStorager(storager.Bolt, path.Join(config.Global.Storage.DataDir, "bolt.db"))
	default:
		db = storager.NewStorager(storager.Bolt, path.Join(config.Global.Storage.DataDir, "bolt.db"))
	}
}

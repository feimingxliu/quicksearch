package core

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/pkg/storager"
	"path"
	"strings"
)

var meta storager.Storager

func InitMeta() {
	switch strings.ToLower(config.Global.Storage.Type) {
	case "bolt":
		meta = storager.NewStorager(storager.Bolt, path.Join(config.Global.Storage.DataDir, "metadata", "bolt.db"))
	default:
		meta = storager.NewStorager(storager.Default, path.Join(config.Global.Storage.DataDir, "metadata", "bolt.db"))
	}
}

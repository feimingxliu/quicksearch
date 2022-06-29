package core

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/pkg/storager"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"log"
	"path"
	"strings"
)

var meta storager.Storager

func InitMeta() {
	switch strings.ToLower(config.Global.Storage.MetaType) {
	case "bolt":
		meta = storager.NewStorager(storager.Bolt, path.Join(config.Global.Storage.DataDir, "metadata", "meta"))
	default:
		meta = storager.NewStorager(storager.Default, path.Join(config.Global.Storage.DataDir, "metadata", "meta"))
	}
	if err := loadAllIndices(); err != nil {
		log.Fatalf("%+v\n", errors.WithStack(err))
	}
}

func loadAllIndices() error {
	indices, err := ListIndices()
	if err != nil {
		return err
	}
	for _, index := range indices {
		if err := index.Open(); err != nil {
			return err
		}
	}
	return nil
}

package jieba

import (
	"embed"
	"github.com/feimingxliu/quicksearch/pkg/util"
	"io"
	"os"
	"path"
)

//go:embed dict/*
var dict embed.FS

var DictDir string

// release all embed files to data dir.
func init() {
	DictDir = path.Join("data", "dict", Name)
	err := os.MkdirAll(DictDir, 0600)
	if err != nil {
		panic(err)
	}
	de, err := dict.ReadDir("dict")
	if err != nil {
		panic(err)
	}
	for _, e := range de {
		p := path.Join(DictDir, e.Name())
		exist, err := util.FileExists(p)
		if err != nil {
			panic(err)
		}
		if !exist {
			outputFile, err := os.Create(p)
			if err != nil {
				panic(err)
			}
			inputFile, err := dict.Open(path.Join("dict", e.Name()))
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(outputFile, inputFile)
			if err != nil {
				panic(err)
			}
		}
	}
}

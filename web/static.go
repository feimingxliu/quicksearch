package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed web/dist/*
var content embed.FS
var StaticFiles fs.FS

func init() {
	var err error
	if StaticFiles, err = fs.Sub(content, "web/dist"); err != nil {
		log.Println("fs.Sub: ", err)
		return
	}
}

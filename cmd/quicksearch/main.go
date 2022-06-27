package main

import (
	"flag"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/core"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/gin"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"log"
)

func main() {
	// init config and index meta.
	InitConfigMeta()
	// run http server.
	if err := gin.ListenAndServe(); err != nil {
		log.Printf("gin.ListenAndServe: %+v", errors.WithStack(err))
	}
	log.Println("close all indies.")
	if err := core.CloseIndices(); err != nil {
		log.Println("core.CloseIndices: ", err)
	}
}

var configPath *string

func init() {
	configPath = flag.String("c", "config.yaml", "the config file path.")
	flag.Parse()
}

func InitConfigMeta() {
	if err := config.Init(*configPath); err != nil {
		log.Fatalln(err)
	}
	core.InitMeta()
}

package main

import (
	"flag"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/core"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/gin"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"log"
	"os"
	"os/signal"
)

func main() {
	// init config.
	InitConfig()
	// start search engine.
	engine := core.NewEngine()
	if err := engine.Run(); err != nil {
		log.Printf("engine.Run: %+v", errors.WithStack(err))
	}
	// run http server.
	if err := gin.ListenAndServe(); err != nil {
		log.Printf("gin.ListenAndServe: %+v", errors.WithStack(err))
	}
	Shutdown()
	// stop engine.
	log.Println("stop engine...")
	if err := engine.Stop(); err != nil {
		log.Printf("engine.Stop: %+v", errors.WithStack(err))
	}
}

var configPath *string

func init() {
	configPath = flag.String("c", "config.yaml", "the config file path.")
	flag.Parse()
}

func InitConfig() {
	if err := config.Init(*configPath); err != nil {
		log.Fatalln(err)
	}
}

func Shutdown() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		// twice ctrl c received
		os.Exit(1)
	}()
}

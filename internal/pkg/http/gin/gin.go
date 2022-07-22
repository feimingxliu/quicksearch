package gin

import (
	"context"
	"fmt"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/middlewares"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/routers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func ListenAndServe() error {
	if config.Global.Env == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middlewares.Cors())
	if config.Global.Env == "debug" {
		engine.Use(gin.Logger())
	}
	routers.RegisterRoutes(engine)
	server := &http.Server{
		Addr:    config.Global.Http.ServerAddr,
		Handler: engine,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("http server Shutdown: %v\n", err)
		}
		log.Println("server shutdown!")
		close(idleConnsClosed)
	}()
	log.Println("server listening at ", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("http server ListenAndServe: %v", err)
	}
	//wait for all idle connection closed.
	<-idleConnsClosed
	return nil
}

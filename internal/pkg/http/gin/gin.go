package gin

import (
	"context"
	"fmt"
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/routers"
	"github.com/feimingxliu/quicksearch/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func ListenAndServe() error {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	if config.Global.Env == "debug" {
		gin.SetMode(gin.DebugMode)
		engine.Use(gin.Logger())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	routers.RegisterRoutes(engine)
	server := &http.Server{
		Addr:    config.Global.Http.ServerAddr,
		Handler: engine,
	}
	staticServer := http.Server{
		Addr:    config.Global.Http.StaticAddr,
		Handler: http.FileServer(http.FS(web.StaticFiles)),
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := staticServer.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("static server Shutdown: %v\n", err)
		}
		log.Println("static server shutdown!")
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("http server Shutdown: %v\n", err)
		}
		log.Println("server shutdown!")
		close(idleConnsClosed)
	}()
	go func() {
		log.Println("static server listening at ", staticServer.Addr)
		if err := staticServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("static server ListenAndServe: %v", err)
		}
	}()
	log.Println("server listening at ", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("http server ListenAndServe: %v", err)
	}
	//wait for all idle connection closed.
	<-idleConnsClosed
	return nil
}

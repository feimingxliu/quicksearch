package routers

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)

func registerPProf(engine *gin.RouterGroup) {
	engine.Any("/", func(c *gin.Context) {
		pprof.Index(c.Writer, c.Request)
	})
	engine.Any("/cmdline", func(c *gin.Context) {
		pprof.Cmdline(c.Writer, c.Request)
	})
	engine.Any("/profile", func(c *gin.Context) {
		pprof.Profile(c.Writer, c.Request)
	})
	engine.Any("/symbol", func(c *gin.Context) {
		pprof.Symbol(c.Writer, c.Request)
	})
	engine.Any("/trace", func(c *gin.Context) {
		pprof.Trace(c.Writer, c.Request)
	})
	engine.Any("/goroutine", func(c *gin.Context) {
		pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
	})
	engine.Any("/threadcreate", func(c *gin.Context) {
		pprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request)
	})
	engine.Any("/heap", func(c *gin.Context) {
		pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
	})
	engine.Any("/allocs", func(c *gin.Context) {
		pprof.Handler("allocs").ServeHTTP(c.Writer, c.Request)
	})
	engine.Any("/block", func(c *gin.Context) {
		pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
	})
	engine.Any("/mutex", func(c *gin.Context) {
		pprof.Handler("mutex").ServeHTTP(c.Writer, c.Request)
	})
}

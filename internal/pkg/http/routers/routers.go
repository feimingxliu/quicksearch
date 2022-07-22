package routers

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/middlewares"
	"github.com/feimingxliu/quicksearch/pkg/about"
	"github.com/feimingxliu/quicksearch/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(engine *gin.Engine) {
	if config.Global.Env == "debug" {
		pprof := engine.Group("/debug/pprof")
		registerPProf(pprof)
	}
	engine.Handle("GET", "/", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/web/")
	})
	static := engine.Group("/web")
	static.StaticFS("/", http.FS(web.StaticFiles))
	if config.Global.Http.Auth.Enabled {
		user := engine.Group("/_user")
		registerUserApi(user)
	}
	v1 := engine.Group("/")
	if config.Global.Http.Auth.Enabled {
		v1.Use(middlewares.Auth())
	}
	{
		//version
		v1.GET("/_version", about.GetVersion)
		index := v1.Group("")
		registerIndexApi(index)
		registerDocumentApi(index)
		registerSearchApi(index)
	}
}

package routers

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/about"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	if config.Global.Env == "debug" {
		pprof := engine.Group("/debug/pprof")
		registerPProf(pprof)
	}
	v1 := engine.Group("/")
	{
		//version
		v1.GET("/_version", about.GetVersion)
		index := v1.Group("")
		registerIndexApi(index)
		registerDocumentApi(index)
		registerSearchApi(index)
	}
}

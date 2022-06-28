package routers

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	if config.Global.Env == "debug" {
		pprof := engine.Group("/debug/pprof")
		registerPProf(pprof)
	}
	v1 := engine.Group("/api/v1")
	{
		index := v1.Group("/index")
		registerIndexApi(index)
		registerDocumentApi(index)
		registerSearchApi(index)
	}
}
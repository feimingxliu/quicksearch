package routers

import (
	"github.com/feimingxliu/quicksearch/internal/pkg/http/handlers/index"
	"github.com/gin-gonic/gin"
)

func registerIndexApi(r *gin.RouterGroup) {
	// create index
	r.POST("/:index", index.Create)
	r.PUT("/:index", index.Create)
	// list indices
	r.GET("/_all", index.List)
}

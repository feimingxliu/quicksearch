package routers

import (
	"github.com/feimingxliu/quicksearch/internal/pkg/http/handlers/index"
	"github.com/gin-gonic/gin"
)

func registerIndexApi(r *gin.RouterGroup) {
	// create index
	r.POST("/:index", index.Create)
	// update index mapping
	r.PUT("/:index/_mapping", index.UpdateMapping)
	// delete index
	r.DELETE("/:index", index.Delete)
	// get index
	r.GET("/:index", index.Get)
	// clone index
	r.POST("/:index/_clone/:target", index.Clone)
	// open index
	r.POST("/:index/_open", index.Open)
	// close index
	r.POST("/:index/_close", index.Close)
	// list indices
	r.GET("/_all", index.List)
}

package routers

import (
	"github.com/feimingxliu/quicksearch/internal/pkg/http/handlers/search"
	"github.com/gin-gonic/gin"
)

func registerSearchApi(r *gin.RouterGroup) {
	// search in index
	r.GET("/:index/_search", search.Search)
	r.POST("/:index/_search", search.Search)
}

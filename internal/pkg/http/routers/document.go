package routers

import (
	"github.com/feimingxliu/quicksearch/internal/pkg/http/handlers/document"
	"github.com/gin-gonic/gin"
)

func registerDocumentApi(r *gin.RouterGroup) {
	// index document
	r.POST("/:index/_doc", document.Index)
	r.POST("/:index/_doc/:id", document.Index)
	// update document
	r.PUT("/:index/_doc/:id", document.Update)
	// bulk index document
	r.POST("/:index/_bulk", document.Bulk)
}

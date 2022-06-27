package document

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/internal/core"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(ctx *gin.Context) {
	indexName := ctx.Param("index")
	if len(indexName) == 0 {
		ctx.JSON(http.StatusBadRequest, "index required!")
		return
	}
	source := make(map[string]interface{})
	if err := ctx.ShouldBindJSON(&source); err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	index, err := core.GetIndex(indexName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	doc := core.NewDocument(source)
	if id := ctx.Param("id"); len(id) > 0 {
		doc.WithID(id)
	}
	err = index.IndexDocument(doc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	ctx.JSON(http.StatusOK, types.IndexDocument{
		Common: types.Common{Acknowledged: true},
		Index:  index.Name,
		ID:     doc.ID,
		Type:   "_doc",
	})
}

func Update(ctx *gin.Context) {

}

func Bulk(ctx *gin.Context) {

}

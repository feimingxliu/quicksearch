package index

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/internal/core"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Create(ctx *gin.Context) {
	indexName := ctx.Param("index")
	if len(indexName) == 0 {
		ctx.JSON(http.StatusBadRequest, "index required!")
		return
	}
	core.NewIndex(core.WithName(indexName))
	ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})
}

func Clone(ctx *gin.Context) {

}

func Open(ctx *gin.Context) {

}

func Close(ctx *gin.Context) {

}

func Delete(ctx *gin.Context) {

}

func Get(ctx *gin.Context) {

}

func List(ctx *gin.Context) {
	indices, err := core.ListIndices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	ctx.JSON(http.StatusOK, types.Indices{
		Common: types.Common{
			Acknowledged: true,
		},
		Indices: indices,
	})
}

func Update(ctx *gin.Context) {
	//TODO
}

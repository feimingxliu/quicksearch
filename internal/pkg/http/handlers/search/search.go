package search

import (
	"github.com/feimingxliu/quicksearch/internal/core"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Search(ctx *gin.Context) {
	searchRequest := new(core.SearchRequest)
	if err := ctx.ShouldBindJSON(searchRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
		return
	}
	var (
		index *core.Index
		res   *core.SearchResult
		err   error
	)
	indexName := ctx.Param("index")
	if len(indexName) > 0 {
		index, err = core.GetIndex(indexName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
			return
		}
	}
	if index == nil {
		res, err = core.Search(searchRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
			return
		}
	} else {
		res, err = index.Search(searchRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, res)
}

package search

import (
	"github.com/gin-gonic/gin"
)

func Search(ctx *gin.Context) {
	/*indexName := ctx.Param("index")
	if len(indexName) == 0 {
		ctx.JSON(http.StatusBadRequest, "index required!")
		return
	}
	params := types.Search{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	if len(params.Query) == 0 {
		ctx.JSON(http.StatusBadRequest, "query required!")
	}
	if params.Timeout < 0 {
		params.Timeout = 0
	}
	if params.TopN < 0 {
		params.TopN = 10
	}
	index, err := core.GetIndex(indexName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	option := &core.SearchOption{}
	result := index.Search(option.SetQuery(params.Query).SetTimeout(time.Duration(params.Timeout) * time.Second).SetTopN(params.TopN))
	ctx.JSON(http.StatusOK, types.SearchResult{
		Common:       types.Common{Acknowledged: true},
		SearchResult: result,
	})*/
}

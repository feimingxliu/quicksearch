package index

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/internal/core"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/types"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func Create(ctx *gin.Context) {
	indexName := ctx.Param("index")
	if len(indexName) == 0 {
		ctx.JSON(http.StatusBadRequest, "index required!")
		return
	}
	body := new(CreateIndex)
	if err := ctx.ShouldBindJSON(body); err != nil {
		if err != io.EOF {
			ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
			return
		}
	}
	options := make([]core.Option, 0)
	if body.Settings != nil {
		options = append(options, core.WithShards(body.Settings.NumberOfShards))
	}
	if body.Mappings != nil {
		options = append(options, core.WithIndexMapping(body.Mappings))
	}
	options = append(options, core.WithName(indexName))
	if _, err := core.NewIndex(options...); err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})
}

func Clone(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	target := ctx.Param("target")
	if len(target) == 0 {
		ctx.JSON(http.StatusBadRequest, "target index required!")
		return
	}
	if err := index.Clone(target); err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})
}

func Open(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	err := index.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})
}

func Close(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	err := index.Close()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})
}

func Delete(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	err := index.Delete()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})
}

func Get(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	ctx.JSON(http.StatusOK, index)
}

func List(ctx *gin.Context) {
	indices, err := core.ListIndices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	ctx.JSON(http.StatusOK, indices)
}

func UpdateMapping(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	mapping := new(core.IndexMapping)
	if err := ctx.ShouldBindJSON(mapping); err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
		return
	}
	_ = index.SetMapping(mapping)
	err := index.UpdateMetadata()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})
}

func getIndex(ctx *gin.Context) (*core.Index, bool) {
	indexName := ctx.Param("index")
	if len(indexName) == 0 {
		ctx.JSON(http.StatusBadRequest, "index required!")
		return nil, false
	}
	index, err := core.GetIndex(indexName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
		return nil, false
	}
	return index, true
}

package index

import (
	"github.com/gin-gonic/gin"
)

func Create(ctx *gin.Context) {
	/*	indexName := ctx.Param("index")
		if len(indexName) == 0 {
			ctx.JSON(http.StatusBadRequest, "index required!")
			return
		}
		core.NewIndex(core.WithName(indexName))
		ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})*/
}

func Clone(ctx *gin.Context) {
	/*	indexName := ctx.Param("index")
		if len(indexName) == 0 {
			ctx.JSON(http.StatusBadRequest, "index required!")
			return
		}
		target := ctx.Param("target")
		if len(target) == 0 {
			ctx.JSON(http.StatusBadRequest, "target index required!")
			return
		}
		if indexName == target {
		ctx.JSON(http.StatusBadRequest, types.Common{Acknowledged: false, Error: errors.ErrCloneIndexSameName.Error()})
		return
		}
		index, err := core.GetIndex(indexName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
			return
		}
		if err := index.Clone(target); err != nil {
			ctx.JSON(http.StatusInternalServerError, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
			return
		}
		ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})*/
}

func Open(ctx *gin.Context) {

}

func Close(ctx *gin.Context) {

}

func Delete(ctx *gin.Context) {

}

func Get(ctx *gin.Context) {
	/*indexName := ctx.Param("index")
	if len(indexName) == 0 {
		ctx.JSON(http.StatusBadRequest, "index required!")
		return
	}
	index, err := core.GetIndex(indexName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	ctx.JSON(http.StatusOK, types.Index{
		Common: types.Common{Acknowledged: true},
		Index:  index,
	})*/
}

func List(ctx *gin.Context) {
	/*indices, err := core.ListIndices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Acknowledged: false, Error: fmt.Sprintf("%+v", err)})
		return
	}
	ctx.JSON(http.StatusOK, types.Indices{
		Common: types.Common{
			Acknowledged: true,
		},
		Indices: indices,
	})*/
}

func Update(ctx *gin.Context) {
	//TODO
}

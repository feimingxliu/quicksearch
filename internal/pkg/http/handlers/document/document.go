package document

import (
	"github.com/feimingxliu/quicksearch/internal/core"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/types"
	"github.com/feimingxliu/quicksearch/pkg/util/uuid"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
)

func Index(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	source := make(map[string]interface{})
	if err := ctx.ShouldBindJSON(&source); err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
		return
	}
	var docID string
	if docID = ctx.Param("id"); docID == "" {
		docID = uuid.GetUUID()
	}
	err := index.IndexOrUpdateDocument(docID, source)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, core.NewBulkActionResult(index.Name, docID, "created", 201, nil, getSeqNo()))
}

func Update(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	var docID string
	if docID = ctx.Param("id"); docID == "" {
		ctx.JSON(http.StatusBadRequest, "doc ID required!")
		return
	}
	fields := make(map[string]interface{})
	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
		return
	}
	err := index.UpdateDocumentPartially(docID, fields)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, core.NewBulkActionResult(index.Name, docID, "updated", 200, nil, getSeqNo()))
}

func Bulk(ctx *gin.Context) {
	body := ctx.Request.Body
	defer body.Close()
	indexName := ctx.Param("index")
	res, err := core.Bulk(indexName, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func Get(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	var docID string
	if docID = ctx.Param("id"); docID == "" {
		ctx.JSON(http.StatusBadRequest, "doc ID required!")
		return
	}
	doc, err := index.GetDocument(docID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, doc)
}

func Delete(ctx *gin.Context) {
	index, ok := getIndex(ctx)
	if !ok {
		return
	}
	var docID string
	if docID = ctx.Param("id"); docID == "" {
		ctx.JSON(http.StatusBadRequest, "doc ID required!")
		return
	}
	err := index.DeleteDocument(docID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, core.NewBulkActionResult(index.Name, docID, "deleted", 200, nil, getSeqNo()))
}

func getIndex(ctx *gin.Context) (*core.Index, bool) {
	indexName := ctx.Param("index")
	if len(indexName) == 0 {
		ctx.JSON(http.StatusBadRequest, "index required!")
		return nil, false
	}
	index, err := core.NewIndex(core.WithName(indexName))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.Common{Error: err.Error()})
		return nil, false
	}
	return index, true
}

var seqNo int64

func getSeqNo() int64 {
	res := seqNo
	atomic.AddInt64(&seqNo, 1)
	return res
}

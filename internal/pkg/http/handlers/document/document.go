package document

import (
	"bufio"
	"fmt"
	"github.com/feimingxliu/quicksearch/internal/core"
	"github.com/feimingxliu/quicksearch/internal/pkg/http/types"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/gin-gonic/gin"
	"log"
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
	index := core.NewIndex(core.WithName(indexName))
	doc := core.NewDocument(source)
	if id := ctx.Param("id"); len(id) > 0 {
		doc.WithID(id)
	}
	err := index.IndexDocument(doc)
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
	indexName := ctx.Param("index")
	if len(indexName) == 0 {
		ctx.JSON(http.StatusBadRequest, "index required!")
		return
	}
	index := core.NewIndex(core.WithName(indexName))
	docs := make([]*core.Document, 0)
	defer ctx.Request.Body.Close()
	scanner := bufio.NewScanner(ctx.Request.Body)
	const maxCapacityPerLine = 1024 * 1024 // 1 MB
	buf := make([]byte, maxCapacityPerLine)
	scanner.Buffer(buf, maxCapacityPerLine)
	for scanner.Scan() {
		source := make(map[string]interface{})
		if err := json.Unmarshal(scanner.Bytes(), &source); err != nil {
			ctx.JSON(http.StatusBadRequest, "invalid line format!")
			return
		}
		docs = append(docs, core.NewDocument(source))
	}
	go func() {
		if err := index.BulkDocuments(docs); err != nil {
			log.Printf("index.BulkDocuments: %+v\n", errors.WithStack(err))
		}
	}()
	ctx.JSON(http.StatusOK, types.Common{Acknowledged: true})
}

func Get(ctx *gin.Context) {

}

func Delete(ctx *gin.Context) {

}

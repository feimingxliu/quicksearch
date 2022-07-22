package routers

import (
	"github.com/feimingxliu/quicksearch/internal/pkg/http/handlers/user"
	"github.com/gin-gonic/gin"
)

func registerUserApi(r *gin.RouterGroup) {
	r.POST("login", user.Login)
}

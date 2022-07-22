package middlewares

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/util/bcrypt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verified := false
		if u, p, ok := ctx.Request.BasicAuth(); ok {
			if u == config.Global.Http.Auth.Username {
				if err := bcrypt.CompareHashAndPassword([]byte(config.Global.Http.Auth.Password), []byte(p)); err == nil {
					verified = true
				}
			}
		}
		if !verified {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

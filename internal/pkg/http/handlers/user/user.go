package user

import (
	"github.com/feimingxliu/quicksearch/internal/config"
	"github.com/feimingxliu/quicksearch/pkg/util/bcrypt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginResult struct {
	Success bool `json:"success"`
}

type LoginForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func Login(ctx *gin.Context) {
	form := new(LoginForm)
	if err := ctx.ShouldBind(form); err != nil {
		ctx.JSON(http.StatusBadRequest, LoginResult{})
		return
	}
	verified := false
	if form.Username == config.Global.Http.Auth.Username {
		if err := bcrypt.CompareHashAndPassword([]byte(config.Global.Http.Auth.Password), []byte(form.Password)); err == nil {
			verified = true
		}
	}
	if !verified {
		ctx.JSON(http.StatusUnauthorized, LoginResult{Success: false})
	} else {
		ctx.JSON(http.StatusOK, LoginResult{Success: true})
	}
}

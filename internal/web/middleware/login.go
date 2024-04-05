package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			//注册和登录是不需要交校验的
			return
		}

		sess := sessions.Default(ctx)
		if sess.Get("userId") == nil {
			//没有话获取到id，中断，不再执行
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

package middleware

import (
	"fmt"
	"net/http"
	"time"

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
		userId := sess.Get("userId")
		if userId == nil {
			//没有话获取到id，中断，不再执行
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		const updateTimeKey = "sess_update_time"
		val := sess.Get(updateTimeKey)

		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Minute*10 {
			//登陆态的时间是15min，最后5min分钟的时候，更新updateTimeKey
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId) //因为sess是一起更新原因，这里要顺便将userId再保存一遍
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

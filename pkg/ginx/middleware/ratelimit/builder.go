package ratelimit

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"example.com/basic-gin/webook/pkg/limiter"
	"github.com/gin-gonic/gin"
)

type Builder struct {
	prefix  string
	limiter limiter.Limiter
}

func NewBuilder(l limiter.Limiter) *Builder {
	return &Builder{
		prefix:  "ip-limiter",
		limiter: l,
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader("x-stress") == "true" {
			//使用context.Context带来的这个标记
			newCtx := context.WithValue(ctx, "x-stress", true)
			ctx.Request = ctx.Request.Clone(newCtx)
			ctx.Next()
			return
		}

		limited, err := b.limiter.Limit(ctx, fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP()))
		if err != nil {
			log.Println(err)

			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if limited { //满足限流条件
			log.Panicln(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

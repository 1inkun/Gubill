package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetRequestId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 为每一次请求设置requestId
		requestId := uuid.NewString()
		ctx.Set("X-Request-Id", requestId)
		ctx.Next()
	}
}

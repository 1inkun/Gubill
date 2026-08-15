package middlewares

import (
	"net/http"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"github.com/gin-gonic/gin"
)

func CheckJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := models.NewResponse()
		// 从请求头中获取JWT
		now := time.Now()
		JWTString := ctx.GetHeader("Authorization")
		if JWTString == "" {
			resp.Code = 403
			resp.Msg = "JWT有误"
			ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}
		// 解析JWT
		claim, err := utils.ParseJWT(JWTString)
		if err != nil {
			resp.Code = 403
			resp.Msg = "JWT有误"
			ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}
		if claim.ExpiresAt.Before(now) {
			resp.Code = 403
			resp.Msg = "JWT已过期"
			ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}
		ctx.Set("Role", claim.Role)
		ctx.Set("userId", claim.UserId)
		ctx.Next()
	}
}

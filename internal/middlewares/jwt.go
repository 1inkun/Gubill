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
		// 从请求头中获取JWT
		now := time.Now()
		JWTString := ctx.GetHeader("Authorization")
		if JWTString == "" {
			resp := models.NewResponse(403, "fail", "登录状态有误", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}
		// 解析JWT
		claim, err := utils.ParseJWT(JWTString)
		if err != nil {
			resp := models.NewResponse(403, "fail", "登录状态有误", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}
		if claim.ExpiresAt.Before(now) {
			resp := models.NewResponse(403, "fail", "登录状态过期", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}
		ctx.Set("Role", claim.Role)
		ctx.Set("userId", claim.UserId)
		// log.Printf("role:%s,userId:%s", claim.Role, claim.UserId)
		ctx.Next()
	}
}

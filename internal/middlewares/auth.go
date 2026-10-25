package middlewares

import (
	"net/http"
	"strings"

	"github.com/1inkun/billApp/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Auth(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := gin.H{
			"code": 0,
			"msg":  "认证失败",
			"data": nil,
		}
		auth := c.GetHeader("Authorization")
		ctx := c.Request.Context()
		if auth == "" {
			c.JSON(http.StatusUnauthorized, resp)
			c.Abort()
			return
		}
		_, tokenStr, ok := strings.Cut(auth, " ")
		if !ok {
			c.JSON(http.StatusUnauthorized, resp)
			c.Abort()
			return
		}
		claims, err := utils.ParseJWT(tokenStr)
		// log.Println(claims.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, resp)
			c.Abort()
			return
		}
		// 验证tokenStr
		check, err := rdb.Get(ctx, claims.Username).Result()
		// log.Printf("Check:%s", check)
		if err != nil {
			c.JSON(http.StatusUnauthorized, resp)
			c.Abort()
			return
		}
		// log.Printf("TokenStr:%s", tokenStr)
		if check != tokenStr {
			c.JSON(http.StatusUnauthorized, resp)
			c.Abort()
			return
		}
		c.Set("role", claims.Role)
		c.Set("userId", claims.UserId)
		c.Next()
	}
}

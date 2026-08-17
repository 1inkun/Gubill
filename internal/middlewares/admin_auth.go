package middlewares

import (
	"net/http"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"github.com/gin-gonic/gin"
)

const (
	// AdminCookieName 管理页面会话 Cookie 名称
	AdminCookieName = "gubill_token"
	// AdminCookieMaxAge Cookie 有效期（秒），与 JWT 有效期保持一致
	AdminCookieMaxAge = 6 * 60 * 60
)

// SetAdminCookie 将 JWT 写入 httpOnly Cookie。
func SetAdminCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AdminCookieName, token, AdminCookieMaxAge, "/", "", false, true)
}

// ClearAdminCookie 清除管理页面会话 Cookie。
func ClearAdminCookie(c *gin.Context) {
	c.SetCookie(AdminCookieName, "", -1, "/", "", false, true)
}

// AdminPageAuth 管理页面鉴权中间件：
// 从 Cookie 读取 JWT，校验有效性并要求 Role=Admin，未登录跳转登录页。
func AdminPageAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(AdminCookieName)
		if err != nil || token == "" {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		claims, err := utils.ParseJWT(token)
		if err != nil || claims.ExpiresAt.Before(time.Now()) {
			ClearAdminCookie(c)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		if claims.Role != "Admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, models.NewResponse(403, "fail", "权限不足", nil))
			return
		}
		c.Set("userId", claims.UserId)
		c.Set("Role", claims.Role)
		c.Next()
	}
}

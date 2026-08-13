package middlewares

import (
	"net/http"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/gin-gonic/gin"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			// log.Panicln("出现错误")
			resp := models.NewResponse()
			err := c.Errors.Last().Err
			// 更新错误响应
			resp["status"] = "fail"
			resp["code"] = 500
			resp["msg"] = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		}
	}
}

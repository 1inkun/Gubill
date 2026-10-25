package middlewares

import (
	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/utils"
	"github.com/gin-gonic/gin"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if be, ok := err.(*models.BusinessErr); ok {
				utils.NewErrResponse(c, be.HttpCode, be.ErrCode, be.Msg, be.Detail)
				return
			} else {
				utils.NewErrResponse(c, 500, 0, "内部错误", "")
				return
			}
		}
	}
}

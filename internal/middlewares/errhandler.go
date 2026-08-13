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
			resp.Status = "fail"
			err := c.Errors.Last().Err
			switch e := err.(type) {
			// 	出现内部错误的情况
			case *models.InternalError:
				{
					resp.Code = e.Code
					resp.Msg = e.Msg
					c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
					break
				}
			case *models.BusinessError:
				{
					resp.Code = e.Code
					resp.Msg = e.Msg
					c.AbortWithStatusJSON(e.Code, resp)
					break
				}
			default:
				{
					resp.Code = 500
					resp.Msg = "内部错误"
					c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
				}
			}
		}
	}
}

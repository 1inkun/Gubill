package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/gin-gonic/gin"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			// log.Panicln("出现错误")
			now := time.Now()
			err := c.Errors.Last().Err
			switch e := err.(type) {
			// 	出现内部错误的情况
			case *models.InternalError:
				{
					resp := models.NewResponse(e.Code, "fail", e.Msg, nil)
					// 内部错误记录在日志中
					log.Printf("%s:%s", now, e.InterErr.Error())
					c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
					break
				}
			case *models.BusinessError:
				{
					resp := models.NewResponse(e.Code, "fail", e.Msg, nil)
					c.AbortWithStatusJSON(e.Code, resp)
					break
				}
			default:
				{
					resp := models.NewResponse(500, "fail", "内部错误", nil)
					c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
				}
			}
		}
	}
}

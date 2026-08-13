package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewResponse() map[string]any {
	var Response = map[string]any{
		"code":   200,
		"status": "success",
		"msg":    "",
		"data":   gin.H{},
	}
	return Response
}

func DataBindErr(c *gin.Context, code int, msg string) {
	resp := NewResponse()
	resp["code"] = code
	resp["msg"] = msg
	c.JSON(http.StatusBadRequest, resp)
}

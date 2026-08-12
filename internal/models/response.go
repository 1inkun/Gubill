package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var Response = map[string]any{
	"code":   200,
	"status": "success",
	"msg":    "",
	"data":   gin.H{},
}

func DataBindErr(c *gin.Context, code int, msg string) {
	resp := Response
	resp["code"] = code
	resp["msg"] = msg
	c.JSON(http.StatusBadRequest, resp)
}

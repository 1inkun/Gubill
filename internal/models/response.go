package models

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
}

func NewResponse() Response {
	var Response = Response{
		200,
		"success",
		"",
		gin.H{},
	}
	return Response
}

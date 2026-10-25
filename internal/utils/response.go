package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一响应
type Response struct {
	Code   int    `josn:"code" example:"200"`
	Status string `json:"status" example:"success"`
	Msg    string `json:"msg" example:"请求成功"`
	Data   any    `json:"data"`
}

type ErrResponse struct {
	ErrCode int    `json:"errCode" example:"0"`
	Status  string `json:"status" example:"fail"`
	Msg     string `json:"msg" example:"内部错误"`
	Detail  any    `json:"detail"`
}

type PaginationData struct {
	Data     any   `json:"data"`
	Page     int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
}

func NewResponse(ctx *gin.Context, code int, msg string, data any) {
	// 此处为正常的逻辑
	requestId, _ := ctx.Get("X-Request-Id")
	ctx.Header("X-Request-Id", requestId.(string))
	resp := Response{
		Code:   0,
		Status: "success",
		Msg:    msg,
		Data:   data,
	}
	ctx.JSON(http.StatusOK, resp)
}

func NewErrResponse(ctx *gin.Context, HttpCode int, ErrCode int, msg string, detail any) {
	requestId, _ := ctx.Get("X-Request-Id")
	ctx.Header("X-Request-Id", requestId.(string))
	resp := ErrResponse{
		ErrCode: ErrCode,
		Status:  "fail",
		Msg:     msg,
		Detail:  detail,
	}
	ctx.JSON(HttpCode, resp)
}

package handler

import (
	"net/http"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type PayHandler struct {
	payService *service.PayService
}

func NewPayHandler(payService *service.PayService) *PayHandler {
	return &PayHandler{payService: payService}
}

func (h *PayHandler) GetUserPayOrders(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	var pagin models.Pagin
	c.MustBindWith(&pagin, binding.Query)
	// 服务层
	res, err := h.payService.GetUserPayOrders(ctx, userId.(string), pagin.Page, pagin.PageSize)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"results":  res,
		"page":     pagin.Page,
		"pageSize": pagin.PageSize,
	}
	c.JSON(http.StatusOK, resp)
}

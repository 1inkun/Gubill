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
	queryData := struct {
		Status int64 `form:"status"`
		models.Pagin
	}{}
	c.MustBindWith(&queryData, binding.Query)
	// 服务层
	res, err := h.payService.GetUserPayOrders(ctx, userId.(string), queryData.Page, queryData.PageSize)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"results":  res,
		"page":     queryData.Page,
		"pageSize": queryData.PageSize,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PayHandler) GetUserPayOrdersById(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	payId := c.Param("payId")
	// 服务层
	res, err := h.payService.GetUserPayOrdersById(ctx, userId.(string), payId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

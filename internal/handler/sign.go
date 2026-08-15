package handler

import (
	"net/http"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
)

type SignHandler struct {
	signService *service.SignService
}

// 面向管理端的接口
func (h *SignHandler) GetAllSignData(c *gin.Context) {
	resp := models.NewResponse()
	resp.Msg = "获取数据成功"
	ctx := c.Request.Context()
	role, exist := c.Get("Role")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	if role != "Admin" {
		c.Error(ErrForbidden)
		return
	}
	pagin := struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}{}
	err := c.ShouldBindQuery(&pagin)
	if err != nil {
		c.Error(err)
		return
	}
	// 服务层处理
	res, err := h.signService.GetAllSignData(ctx, pagin.Page, pagin.PageSize)
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

func (h *SignHandler) GetSignDataBySignId(c *gin.Context) {
	resp := models.NewResponse()
	resp.Msg = "获取数据成功"
	ctx := c.Request.Context()
	role, exist := c.Get("Role")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	if role != "Admin" {
		c.Error(ErrForbidden)
		return
	}
	signId := c.Param("sign_id")
	// 服务层
	res, err := h.signService.GetSignDataBySignId(ctx, signId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

func (h *SignHandler) UpdateSignData(c *gin.Context) {
	resp := models.NewResponse()
	resp.Msg = "更新成功"
	ctx := c.Request.Context()
	role, exist := c.Get("Role")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	if role != "Admin" {
		c.Error(ErrForbidden)
		return
	}
	signId := c.Param("sign_id")
	bodyData := struct {
		Status int64 `json:"status"`
		Value  int64 `json:"value"`
	}{}
	c.ShouldBindJSON(&bodyData)
	// 服务层
	res, err := h.signService.UpdateSignData(ctx, signId, bodyData.Status, bodyData.Value)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"rowsAffected": res,
	}
	c.JSON(http.StatusOK, resp)
}

// 面向用户端的接口
func NewSignHandler(signService *service.SignService) *SignHandler {
	return &SignHandler{signService: signService}
}

func (h *SignHandler) GenerateNewSignData(c *gin.Context) {
	response := models.NewResponse()
	response.Msg = "签到成功"
	// 从上下文中获取用户ID
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	res, err := h.signService.GenerateNewSignData(ctx, userId.(string))
	if err != nil {
		c.Error(err)
		return
	}
	response.Data = gin.H{
		"signId": res,
	}
	c.JSON(http.StatusOK, response)
}

func (h *SignHandler) GetUserSignData(c *gin.Context) {
	resp := models.NewResponse()
	resp.Msg = "获取成功"
	// 从上下文中获取用户ID
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	queryData := struct {
		Status int64 `form:"status"`
	}{}
	err := c.ShouldBindQuery(&queryData)
	if err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 数据处理完毕交由服务层函数处理
	res, err := h.signService.GetUserSignData(ctx, userId.(string), queryData.Status)
	if err != nil {
		c.Error(err)
		return
	}
	if len(res) != 0 {
		resp.Data = res
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SignHandler) GetSignData(c *gin.Context) {
	resp := models.NewResponse()
	resp.Msg = "获取成功"
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	signId := c.Param("sign_id")
	if signId == "" {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层处理
	res, err := h.signService.GetSignData(ctx, userId.(string), signId)
	if err != nil {
		c.Error(err)
		return
	}
	if res.UUID != "" {
		resp.Data = res
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SignHandler) FinishSignData(c *gin.Context) {
	resp := models.NewResponse()
	resp.Data = "结算成功"
	ctx := c.Request.Context()
	signId := c.Param("sign_id")
	if signId == "" {
		c.Error(ErrBindDataError)
		return
	}
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层处理
	res, err := h.signService.FinishSignData(ctx, userId.(string), signId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"pay_id": res,
	}
	c.JSON(http.StatusOK, resp)
}

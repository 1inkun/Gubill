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

func NewSignHandler(signService *service.SignService) *SignHandler {
	return &SignHandler{signService: signService}
}

var (
	ErrNeedLogin = models.NewBusinessError(401, "需要登录")
)

// 面向管理端的接口
func (h *SignHandler) GetAllSignData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取数据成功", nil)
	ctx := c.Request.Context()
	role, _ := c.Get("Role")
	if role != "Admin" {
		c.Error(ErrForbidden)
		return
	}
	// 配置查询参数
	pagin := models.Pagin{}
	err := c.ShouldBindQuery(&pagin)
	if err != nil {
		c.Error(ErrBindDataError)
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
	resp := models.NewResponse(200, "success", "获取数据成功", nil)
	ctx := c.Request.Context()
	role, _ := c.Get("Role")
	if role != "Admin" {
		c.Error(ErrForbidden)
		return
	}
	signId := c.Param("sign_id")
	// 服务层
	res, err := h.signService.GetSignDataById(ctx, signId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

func (h *SignHandler) UpdateSignData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "更新成功", nil)
	ctx := c.Request.Context()
	role, _ := c.Get("Role")
	if role != "Admin" {
		c.Error(ErrForbidden)
		return
	}
	signId := c.Param("sign_id")
	bodyData := struct {
		Status int64 `json:"status" binding:"required"`
		Value  int64 `json:"value" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&bodyData); err != nil {
		c.Error(ErrBindDataError)
		return
	}
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
func (h *SignHandler) GenerateNewSignData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "签到成功", nil)
	// 从上下文中获取用户ID
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrNeedLogin)
		return
	}
	res, err := h.signService.GenerateNewSignData(ctx, userId.(string))
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"signId": res,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SignHandler) GetUserSignData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
	// 从上下文中获取用户ID
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrNeedLogin)
		return
	}
	// 设置查询参数
	queryData := struct {
		Status int64 `form:"status"`
		models.Pagin
	}{}
	err := c.ShouldBindQuery(&queryData)
	if err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 数据处理完毕交由服务层函数处理
	res, err := h.signService.GetUserSignData(ctx, userId.(string), queryData.Status, queryData.Page, queryData.PageSize)
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
	resp := models.NewResponse(200, "success", "获取成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrNeedLogin)
		return
	}
	signId := c.Param("sign_id")
	// 服务层处理
	res, err := h.signService.GetUserSignDataById(ctx, userId.(string), signId)
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
	resp := models.NewResponse(200, "success", "结算成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrNeedLogin)
		return
	}
	signId := c.Param("sign_id")
	// 服务层处理
	res, err := h.signService.FinishSignData(ctx, userId.(string), signId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"payId": res,
	}
	c.JSON(http.StatusOK, resp)
}

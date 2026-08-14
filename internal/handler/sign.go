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

package handler

import (
	"net/http"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService *service.MemberService
}

func NewMemberHandler(memberService *service.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

var (
	ErrForbidden = models.NewBusinessError(403, "权限不足")
)

func (h *MemberHandler) AddNewMemberPlan(c *gin.Context) {
	resp := models.NewResponse()
	resp.Msg = "添加成功"
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
	bodyData := struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value int64  `json:"value"`
		Des   string `json:"des"`
	}{}
	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层处理
	res, err := h.memberService.AddNewMemberPlan(ctx, bodyData.Name, bodyData.Type, bodyData.Value, bodyData.Des)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"planId": res,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MemberHandler) GetMemberPlanData(c *gin.Context) {
	resp := models.NewResponse()
	resp.Msg = "数据获取成功"
	c.JSON(http.StatusOK, resp)
}

package handlers

import (
	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/services"
	"github.com/1inkun/billApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService *services.MemberService
}

func NewMemberHandler(memberService *services.MemberService) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
	}
}

// @Summary					查询用户的会员信息
// @Tags						member
// @Accept						json
// @Produce					json
// @Param						SearchMemberParams	query	models.SearchMemberParams	true	"查询会员信息时需要提供的信息"
// @securityDefinitions.apikey	ApiKeyAuth
// @Success					200	{object}	utils.Response{data=models.MemberData}
// @Failure					400	{object}	utils.ErrResponse{detail=string}
// @Failure					401	{object}	utils.ErrResponse{detail=string}	"登录状态有误"
// @Failure					500	{object}	utils.ErrResponse{detail=nil}
// @Router						/member [get]
func (h *MemberHandler) GetMemberData(c *gin.Context) {
	ctx := c.Request.Context()
	requestId, ok := c.Get("userId")
	if !ok {
		c.Error(models.ErrLoginStatus)
		return
	}
	params := models.SearchMemberParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.Error(models.ErrBindingData)
		return
	}
	if params.UserId != requestId.(int64) {
		c.Error(models.ErrIllegalOperation)
		return
	}
	data, err := h.memberService.GetMemberData(ctx, params.UserId)
	if err != nil {
		c.Error(err)
		return
	}
	utils.NewResponse(c, 200, "获取成功", data)
}

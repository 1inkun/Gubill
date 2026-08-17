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

// 管理侧接口
// MemberPlan相关接口
func (h *MemberHandler) AddNewMemberPlan(c *gin.Context) {
	resp := models.NewResponse(200, "success", "添加成功", nil)
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
		Name  string `json:"name" binding:"required"`
		Type  string `json:"type" binding:"required"`
		Value int64  `json:"value" binding:"required"`
		Des   string `json:"des" binding:"required"`
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

func (h *MemberHandler) UpdateMemberPlan(c *gin.Context) {
	resp := models.NewResponse(200, "success", "修改成功", nil)
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
	planId := c.Param("plan_id")
	bodyData := struct {
		Name  string `json:"name" binding:"required"`
		Value int64  `json:"value" binding:"required"`
		Des   string `json:"des" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&bodyData); err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层
	res, err := h.memberService.UpdatePlanData(ctx, planId, bodyData.Name, bodyData.Value, bodyData.Des)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"rowsAffected": res,
	}
	c.JSON(http.StatusOK, resp)
}

// MemberOrder相关接口
func (h *MemberHandler) GetAllMemberOrderData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
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
		Page     int `form:"page"`
		PageSize int `form:"page_size"`
	}{}
	if err := c.ShouldBindQuery(&pagin); err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层
	res, err := h.memberService.GetAllMemberOrderData(ctx, pagin.Page, pagin.PageSize)
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

func (h *MemberHandler) GetMemberOrderDataByOrderId(c *gin.Context) {
	resp := models.NewResponse(200, "success", "修改成功", nil)
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
	orderId := c.Param("order_id")
	// 服务层
	res, err := h.memberService.GetMemberOrderDataByOrderId(ctx, orderId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

func (h *MemberHandler) UpdateMemberOrderData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "修改成功", nil)
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
	orderId := c.Param("order_id")
	bodyData := struct {
		Value  int64 `json:"value" binding:"required"`
		Status int64 `json:"status" binding:"required"`
	}{}
	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层
	res, err := h.memberService.UpdateMemberOrderData(ctx, orderId, bodyData.Status, bodyData.Value)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"rowsAffected": res,
	}
	c.JSON(http.StatusOK, resp)
}

// MemberList相关接口
func (h *MemberHandler) GetAllMemberListData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
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
		Page     int `form:"page"`
		PageSize int `form:"page_size"`
	}{}
	if err := c.ShouldBindQuery(&pagin); err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层
	res, err := h.memberService.GetAllMemberListData(ctx, pagin.Page, pagin.PageSize)
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

func (h *MemberHandler) GetMemberListDataByMemberId(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
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
	memberId := c.Param("member_id")
	// 服务层
	res, err := h.memberService.GetMemberListDataByMemberId(ctx, memberId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

func (h *MemberHandler) AddNewMemberListData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "添加新的会员信息成功", nil)
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
		UserId string `json:"userId" binding:"required"`
		EndAt  int64  `json:"end_at" binding:"required"`
	}{}
	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层
	res, err := h.memberService.AddNewMemberListData(ctx, bodyData.UserId, bodyData.EndAt)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"memberId": res,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MemberHandler) UpdateMemberListData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "更新会员信息成功", nil)
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
	memberId := c.Param("member_id")
	bodyData := struct {
		Status int64 `json:"status" binding:"required"`
		EndAt  int64 `json:"end_at" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&bodyData); err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层
	res, err := h.memberService.UpdateMemberListData(ctx, memberId, bodyData.Status, bodyData.EndAt)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"rowsAffected": res,
	}
	c.JSON(http.StatusOK, resp)
}

// 用户侧接口
// MemberPlan相关接口
func (h *MemberHandler) GetMemberPlanData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "数据获取成功", nil)
	ctx := c.Request.Context()
	planId := c.Param("plan_id")
	// 服务层处理
	res, err := h.memberService.GetMemberPlanData(ctx, planId)
	if err != nil {
		c.Error(err)
		return
	}
	if res.UUID != "" {
		resp.Data = res
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MemberHandler) GenMemberPlanOrder(c *gin.Context) {
	resp := models.NewResponse(200, "success", "订单生成成功", nil)
	ctx := c.Request.Context()
	planId := c.Param("plan_id")
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层处理
	res, err := h.memberService.GenMemberPlanOrder(ctx, userId.(string), planId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"payId": res,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MemberHandler) GetAllMemberPlans(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取会员计划成功", nil)
	ctx := c.Request.Context()
	// 服务层处理
	res, err := h.memberService.GetAllMemberPlans(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	if res != nil {
		resp.Data = res
	}
	c.JSON(http.StatusOK, resp)
}

// MemberOrder相关接口
func (h *MemberHandler) GetUserMemberOrderData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取数据成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	// 数据量较大,可能要分页查询
	pagin := struct {
		Page     int `form:"page"`
		PageSize int `form:"page_size"`
	}{}
	err := c.ShouldBindQuery(&pagin)
	if err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层
	res, err := h.memberService.GetUserMemberOrderData(ctx, userId.(string), pagin.Page, pagin.PageSize)
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

func (h *MemberHandler) GetUserMemberOrderDataById(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取数据成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	orderId := c.Param("order_id")
	// 服务层
	res, err := h.memberService.GetUserMemberOrderDataById(ctx, userId.(string), orderId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

func (h *MemberHandler) CancelMemberOrder(c *gin.Context) {
	resp := models.NewResponse(200, "success", "操作成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	orderId := c.Param("order_id")
	// 服务层
	res, err := h.memberService.CancelMemberOrder(ctx, userId.(string), orderId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = gin.H{
		"rowsAffected": res,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MemberHandler) FinishMemberOrder(c *gin.Context) {
	resp := models.NewResponse(200, "success", "结算成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	orderId := c.Param("order_id")
	// 服务层
	res, err := h.memberService.FinishMemberOrder(ctx, userId.(string), orderId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

// MemberList相关接口
func (h *MemberHandler) GetUserMemberData(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取会员信息成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	// 服务层
	res, err := h.memberService.GetUserMemberData(ctx, userId.(string))
	if err != nil {
		c.Error(err)
		return
	}
	if res.UUID == "" {
		resp.Msg = "还未购买会员"
	} else {
		resp.Data = res
	}
	c.JSON(http.StatusOK, resp)
}

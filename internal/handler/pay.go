package handler

import (
	"net/http"
	"strconv"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
)

// PayHandler 支付单相关接口
type PayHandler struct {
	payService *service.PaymentService
}

func NewPayHandler(payService *service.PaymentService) *PayHandler {
	return &PayHandler{payService: payService}
}

// 用户侧接口

// GetUserPays 获取用户自己的支付记录（分页）
func (h *PayHandler) GetUserPays(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	pagin := parsePageQuery(c)
	if pagin == nil {
		c.Error(ErrBindDataError)
		return
	}
	res, err := h.payService.ListUserPays(ctx, userId.(string), pagin.Page, pagin.PageSize)
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

// GetUserPay 获取用户自己的单笔支付单
func (h *PayHandler) GetUserPay(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
	ctx := c.Request.Context()
	userId, exist := c.Get("userId")
	if !exist {
		c.Error(ErrBindDataError)
		return
	}
	payId := c.Param("pay_id")
	res, err := h.payService.GetUserPay(ctx, userId.(string), payId)
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

// 管理侧接口

// GetAllPays 获取全部支付记录（分页）
func (h *PayHandler) GetAllPays(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
	ctx := c.Request.Context()
	if !isAdmin(c) {
		c.Error(ErrForbidden)
		return
	}
	pagin := parsePageQuery(c)
	if pagin == nil {
		c.Error(ErrBindDataError)
		return
	}
	res, err := h.payService.ListAllPays(ctx, pagin.Page, pagin.PageSize)
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

// GetPay 获取单笔支付单
func (h *PayHandler) GetPay(c *gin.Context) {
	resp := models.NewResponse(200, "success", "获取成功", nil)
	ctx := c.Request.Context()
	if !isAdmin(c) {
		c.Error(ErrForbidden)
		return
	}
	res, err := h.payService.GetPay(ctx, c.Param("pay_id"))
	if err != nil {
		c.Error(err)
		return
	}
	resp.Data = res
	c.JSON(http.StatusOK, resp)
}

// ConfirmPay 管理员兜底确认收款
func (h *PayHandler) ConfirmPay(c *gin.Context) {
	resp := models.NewResponse(200, "success", "确认收款成功", nil)
	ctx := c.Request.Context()
	if !isAdmin(c) {
		c.Error(ErrForbidden)
		return
	}
	if err := h.payService.ConfirmPaid(ctx, c.Param("pay_id")); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// RefundPay 管理员退款（全额）
func (h *PayHandler) RefundPay(c *gin.Context) {
	resp := models.NewResponse(200, "success", "退款成功", nil)
	ctx := c.Request.Context()
	if !isAdmin(c) {
		c.Error(ErrForbidden)
		return
	}
	if err := h.payService.RefundPay(ctx, c.Param("pay_id")); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

type pageQuery struct {
	Page     int
	PageSize int
}

// parsePageQuery 解析可选分页参数（与 utils.Paginate 语义一致）：
// page 默认/归一为 1，page_size 默认 10、上限 100；仅非数字参数视为参数错误。
func parsePageQuery(c *gin.Context) *pageQuery {
	q := &pageQuery{Page: 1, PageSize: 10}
	if v := c.Query("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}
		if n <= 0 {
			n = 1
		}
		q.Page = n
	}
	if v := c.Query("page_size"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}
		if n <= 0 {
			n = 10
		}
		if n > 100 {
			n = 100
		}
		q.PageSize = n
	}
	return q
}

func isAdmin(c *gin.Context) bool {
	role, exist := c.Get("Role")
	return exist && role == "Admin"
}

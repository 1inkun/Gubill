package handler

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/1inkun/Gubill/internal/middlewares"
	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*.html
var templatesFS embed.FS

// WebHandler 管理页面与站内支付页
type WebHandler struct {
	userService   *service.UserService
	memberService *service.MemberService
	signService   *service.SignService
	payService    *service.PaymentService
	tmpl          *template.Template
}

func NewWebHandler(userService *service.UserService, memberService *service.MemberService,
	signService *service.SignService, payService *service.PaymentService) (*WebHandler, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"yuan": func(v int64) float64 { return float64(v) / 100 },
		"ftime": func(v int64) string {
			if v == 0 {
				return "-"
			}
			return time.Unix(v, 0).Format("2006-01-02 15:04:05")
		},
	}).ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &WebHandler{
		userService:   userService,
		memberService: memberService,
		signService:   signService,
		payService:    payService,
		tmpl:          tmpl,
	}, nil
}

func (h *WebHandler) render(c *gin.Context, name string, data any) {
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(c.Writer, name, data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) redirect(c *gin.Context, path, msg string) {
	if msg != "" {
		c.Redirect(http.StatusFound, path+"?msg="+msg)
		return
	}
	c.Redirect(http.StatusFound, path)
}

// 登录

func (h *WebHandler) LoginPage(c *gin.Context) {
	data := gin.H{"Title": "管理员登录", "Msg": c.Query("msg")}
	h.render(c, "admin_login", data)
}

// HomePage 用户端根路径欢迎页
func (h *WebHandler) HomePage(c *gin.Context) {
	h.render(c, "user_home", gin.H{"Title": "Gubill 用户端"})
}

func (h *WebHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	ctx := c.Request.Context()
	token, err := h.userService.Login(ctx, username, password, c.ClientIP())
	if err != nil {
		data := gin.H{"Title": "管理员登录", "Msg": "用户名或密码错误"}
		c.Status(http.StatusOK)
		h.render(c, "admin_login", data)
		return
	}
	middlewares.SetAdminCookie(c, token)
	h.redirect(c, "/admin", "")
}

func (h *WebHandler) Logout(c *gin.Context) {
	middlewares.ClearAdminCookie(c)
	h.redirect(c, "/admin/login", "")
}

// 概览

func (h *WebHandler) Dashboard(c *gin.Context) {
	ctx := c.Request.Context()
	paidCount, paidSum, err := h.payService.TodayStats(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	pending, err := h.payService.PendingCount(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	plans, err := h.memberService.GetAllMemberPlans(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	data := gin.H{
		"Title":       "控制台",
		"PaidCount":   paidCount,
		"PaidSumYuan": float64(paidSum) / 100,
		"Pending":     pending,
		"PlanCount":   len(plans),
		"Msg":         c.Query("msg"),
	}
	h.render(c, "admin_dashboard", data)
}

// 会员计划

func (h *WebHandler) PlansPage(c *gin.Context) {
	ctx := c.Request.Context()
	plans, err := h.memberService.GetAllMemberPlans(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	data := gin.H{"Title": "会员计划", "Plans": plans, "Msg": c.Query("msg")}
	h.render(c, "admin_plans", data)
}

func (h *WebHandler) CreatePlan(c *gin.Context) {
	ctx := c.Request.Context()
	value, _ := strconv.ParseInt(c.PostForm("value"), 10, 64)
	if _, err := h.memberService.AddNewMemberPlan(ctx,
		c.PostForm("name"), c.PostForm("type"), value, c.PostForm("des")); err != nil {
		h.redirect(c, "/admin/plans", err.Error())
		return
	}
	h.redirect(c, "/admin/plans", "添加成功")
}

func (h *WebHandler) UpdatePlan(c *gin.Context) {
	ctx := c.Request.Context()
	value, _ := strconv.ParseInt(c.PostForm("value"), 10, 64)
	if _, err := h.memberService.UpdatePlanData(ctx,
		c.Param("plan_id"), c.PostForm("name"), value, c.PostForm("des")); err != nil {
		h.redirect(c, "/admin/plans", err.Error())
		return
	}
	h.redirect(c, "/admin/plans", "修改成功")
}

// 会员订单

func (h *WebHandler) OrdersPage(c *gin.Context) {
	ctx := c.Request.Context()
	page, pageSize := pageFromQuery(c)
	orders, err := h.memberService.GetAllMemberOrderData(ctx, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	type orderView struct {
		models.MemberOrderRes
		StatusText string
	}
	views := make([]orderView, 0, len(orders))
	for _, o := range orders {
		views = append(views, orderView{o, models.MemberOrderStatusText(o.Status)})
	}
	data := gin.H{"Title": "会员订单", "Orders": views, "Msg": c.Query("msg")}
	h.render(c, "admin_orders", data)
}

func (h *WebHandler) ConfirmOrder(c *gin.Context) {
	ctx := c.Request.Context()
	orderId := c.Param("order_id")
	pay, err := h.payService.GetPayByBusinessId(ctx, orderId)
	if err != nil {
		h.redirect(c, "/admin/orders", "未找到该订单的支付单")
		return
	}
	if err := h.payService.ConfirmPaid(ctx, pay.UUID); err != nil {
		h.redirect(c, "/admin/orders", err.Error())
		return
	}
	h.redirect(c, "/admin/orders", "确认收款成功")
}

func (h *WebHandler) CancelOrder(c *gin.Context) {
	ctx := c.Request.Context()
	orderId := c.Param("order_id")
	order, err := h.memberService.GetMemberOrderDataByOrderId(ctx, orderId)
	if err != nil {
		h.redirect(c, "/admin/orders", err.Error())
		return
	}
	if _, err := h.memberService.CancelMemberOrder(ctx, order.UserId, orderId); err != nil {
		h.redirect(c, "/admin/orders", err.Error())
		return
	}
	h.redirect(c, "/admin/orders", "取消成功")
}

// 会员管理

func (h *WebHandler) MembersPage(c *gin.Context) {
	ctx := c.Request.Context()
	page, pageSize := pageFromQuery(c)
	members, err := h.memberService.GetAllMemberListData(ctx, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	type memberView struct {
		models.MemberListRes
		StatusText string
	}
	views := make([]memberView, 0, len(members))
	for _, m := range members {
		text := "已失效"
		if m.Status == models.MemberStatusValid {
			text = "有效"
		}
		views = append(views, memberView{m, text})
	}
	data := gin.H{"Title": "会员管理", "Members": views, "Msg": c.Query("msg")}
	h.render(c, "admin_members", data)
}

func (h *WebHandler) AddMember(c *gin.Context) {
	ctx := c.Request.Context()
	endAt, _ := strconv.ParseInt(c.PostForm("end_at"), 10, 64)
	if _, err := h.memberService.AddNewMemberListData(ctx, c.PostForm("userId"), endAt); err != nil {
		h.redirect(c, "/admin/members", err.Error())
		return
	}
	h.redirect(c, "/admin/members", "添加成功")
}

func (h *WebHandler) UpdateMember(c *gin.Context) {
	ctx := c.Request.Context()
	status, _ := strconv.ParseInt(c.PostForm("status"), 10, 64)
	endAt, _ := strconv.ParseInt(c.PostForm("end_at"), 10, 64)
	if _, err := h.memberService.UpdateMemberListData(ctx, c.Param("member_id"), status, endAt); err != nil {
		h.redirect(c, "/admin/members", err.Error())
		return
	}
	h.redirect(c, "/admin/members", "修改成功")
}

// 签到记录

func (h *WebHandler) SignsPage(c *gin.Context) {
	ctx := c.Request.Context()
	page, pageSize := pageFromQuery(c)
	signs, err := h.signService.GetAllSignData(ctx, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	type signView struct {
		models.SignRes
		StatusText string
	}
	views := make([]signView, 0, len(signs))
	for _, s := range signs {
		views = append(views, signView{s, models.SignStatusText(s.Status)})
	}
	data := gin.H{"Title": "签到记录", "Signs": views, "Msg": c.Query("msg")}
	h.render(c, "admin_signs", data)
}

// 支付记录

func (h *WebHandler) PaysPage(c *gin.Context) {
	ctx := c.Request.Context()
	page, pageSize := pageFromQuery(c)
	pays, err := h.payService.ListAllPays(ctx, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	type payView struct {
		models.PayRes
		StatusText string
		AmountYuan float64
	}
	views := make([]payView, 0, len(pays))
	for _, p := range pays {
		views = append(views, payView{p, models.PayStatusText(p.Status), float64(p.Value) / 100})
	}
	data := gin.H{"Title": "支付记录", "Pays": views, "Msg": c.Query("msg")}
	h.render(c, "admin_pays", data)
}

func (h *WebHandler) ConfirmPayPage(c *gin.Context) {
	ctx := c.Request.Context()
	if err := h.payService.ConfirmPaid(ctx, c.Param("pay_id")); err != nil {
		h.redirect(c, "/admin/pays", err.Error())
		return
	}
	h.redirect(c, "/admin/pays", "确认收款成功")
}

func (h *WebHandler) RefundPayPage(c *gin.Context) {
	ctx := c.Request.Context()
	if err := h.payService.RefundPay(ctx, c.Param("pay_id")); err != nil {
		h.redirect(c, "/admin/pays", err.Error())
		return
	}
	h.redirect(c, "/admin/pays", "退款成功")
}

func pageFromQuery(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10
	if v, err := strconv.Atoi(c.Query("page")); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.Atoi(c.Query("page_size")); err == nil && v > 0 {
		pageSize = v
		if pageSize > 100 {
			pageSize = 100
		}
	}
	return page, pageSize
}

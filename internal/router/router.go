package router

import (
	"os"

	"github.com/1inkun/Gubill/internal/handler"
	"github.com/1inkun/Gubill/internal/middlewares"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 面向用户开放的接口
func InitRouter(db *gorm.DB, paymentService *service.PaymentService) *gin.Engine {
	ginMode := os.Getenv("GinMode")
	if ginMode == "" {
		gin.SetMode("debug")
	} else {
		gin.SetMode(ginMode)
	}
	r := gin.Default()
	// 中间件
	r.Use(middlewares.RateLimiter(), middlewares.ErrHandler())

	userService := service.NewUserService(db)
	signService := service.NewSignService(db)
	memberService := service.NewMemberService(db)

	userHandler := handler.NewUserHandler(userService)
	sighHandler := handler.NewSignHandler(signService)
	memberHandler := handler.NewMemberHandler(memberService)
	payHandler := handler.NewPayHandler(paymentService)

	apiv1 := r.Group("/api/v1")
	user := apiv1.Group("/user")
	{
		user.POST("/login", userHandler.Login)
		user.POST("/register", userHandler.Register)
	}
	sign := apiv1.Group("/sign")
	{
		sign.Use(middlewares.CheckJWT())
		sign.POST("", sighHandler.GenerateNewSignData)
		sign.GET("", sighHandler.GetUserSignData)
		sign.GET("/:sign_id", sighHandler.GetSignData)
		sign.PUT("/:sign_id", sighHandler.FinishSignData)
	}
	memberPlan := apiv1.Group("/member_plan")
	{
		memberPlan.Use(middlewares.CheckJWT())
		// 获取全部会员计划
		memberPlan.GET("", memberHandler.GetAllMemberPlans)
		// 获取单个会员计划
		memberPlan.GET("/:plan_id", memberHandler.GetMemberPlanData)
		// 通过会员计划生成会员订单
		memberPlan.POST("/:plan_id", memberHandler.GenMemberPlanOrder)
	}
	memberList := apiv1.Group("/member_list")
	{
		memberList.Use(middlewares.CheckJWT())
		// 获取用户自己的**会员情况**
		memberList.GET("", memberHandler.GetUserMemberData)

	}
	memberOrder := apiv1.Group("/member_order")
	{
		memberOrder.Use(middlewares.CheckJWT())
		// 获取用户自己的**会员订单**
		memberOrder.GET("", memberHandler.GetUserMemberOrderData)
		// 通过订单ID获取用户自己的**会员订单**
		memberOrder.GET("/:order_id", memberHandler.GetUserMemberOrderDataById)
		// 取消用户自己的**会员订单**
		memberOrder.DELETE("/:order_id", memberHandler.CancelMemberOrder)
		// 结算用户自己的**会员订单**
		memberOrder.POST("/:order_id", memberHandler.FinishMemberOrder)
	}
	pay := apiv1.Group("/pay")
	{
		authed := pay.Group("")
		authed.Use(middlewares.CheckJWT())
		authed.GET("", payHandler.GetUserPays)
		authed.GET("/:pay_id", payHandler.GetUserPay)
	}
	return r
}

// 面向管理员开放的接口
func InitAdminRouter(db *gorm.DB, paymentService *service.PaymentService) *gin.Engine {
	ginMode := os.Getenv("GinMode")
	if ginMode == "" {
		gin.SetMode("debug")
	} else {
		gin.SetMode(ginMode)
	}
	r := gin.Default()
	r.Use(middlewares.RateLimiter(), middlewares.ErrHandler())
	// 构造Service
	memberService := service.NewMemberService(db)
	signService := service.NewSignService(db)
	// 构造Handler
	memberHandler := handler.NewMemberHandler(memberService)
	signHandler := handler.NewSignHandler(signService)
	payHandler := handler.NewPayHandler(paymentService)

	apiv1Admin := r.Group("/api/v1/", middlewares.CheckJWT())
	// 会员服务相关的管理接口
	memberPlan := apiv1Admin.Group("/member_plan")
	{
		// 新增会员计划
		memberPlan.POST("", memberHandler.AddNewMemberPlan)
		// 修改现有的会员计划
		memberPlan.PUT("/:plan_id", memberHandler.UpdateMemberPlan)
	}
	memberList := apiv1Admin.Group("/member_list")
	{
		// 获取目前所有用户的会员情况
		memberList.GET("", memberHandler.GetAllMemberListData)
		// 新增会员情况
		memberList.POST("", memberHandler.AddNewMemberListData)
		// 修改会员情况
		memberList.PUT("/:member_id", memberHandler.UpdateMemberListData)
		memberList.GET("/:member_id", memberHandler.GetMemberListDataByMemberId)
	}
	memberOrder := apiv1Admin.Group("/member_order")
	{
		memberOrder.GET("", memberHandler.GetAllMemberOrderData)
		memberOrder.GET("/:order_id", memberHandler.GetMemberOrderDataByOrderId)
		memberOrder.PUT("/:order_id", memberHandler.UpdateMemberOrderData)
	}
	// 签到服务的相关管理接口
	sign := apiv1Admin.Group("/sign")
	{
		sign.GET("", signHandler.GetAllSignData)
		sign.GET("/:sign_id", signHandler.GetSignDataBySignId)
		sign.PUT("/:sign_id", signHandler.UpdateSignData)
	}
	pay := apiv1Admin.Group("/pay")
	{
		pay.GET("", payHandler.GetAllPays)
		pay.GET("/:pay_id", payHandler.GetPay)
		pay.POST("/:pay_id/confirm", payHandler.ConfirmPay)
		pay.POST("/:pay_id/refund", payHandler.RefundPay)
	}
	return r
}

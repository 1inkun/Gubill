package router

import (
	"os"

	"github.com/1inkun/Gubill/internal/handler"
	"github.com/1inkun/Gubill/internal/middlewares"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
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
		memberPlan.GET("", memberHandler.GetAllMemberPlans)
		memberPlan.GET("/:plan_id", memberHandler.AddNewMemberPlan)
		memberPlan.POST("/:plan_id", memberHandler.GenMemberPlanOrder)
	}
	memberList := apiv1.Group("/member_list")
	{
		memberList.Use(middlewares.CheckJWT())
		memberList.GET("", memberHandler.GetUserMemberData)

	}
	memberOrder := apiv1.Group("/member_order")
	{
		memberOrder.Use(middlewares.CheckJWT())
		memberOrder.GET("", memberHandler.GetUserMemberOrderData)
		memberOrder.GET("/:order_id", memberHandler.GetUserMemberOrderDataById)
		memberOrder.DELETE("/:order_id", memberHandler.CancelMemberOrder)
		memberOrder.POST("/:order_id", memberHandler.FinishMemberOrder)
	}
	return r
}

func InitAdminRouter(db *gorm.DB) *gin.Engine {
	ginMode := os.Getenv("GinMode")
	if ginMode == "" {
		gin.SetMode("debug")
	} else {
		gin.SetMode(ginMode)
	}
	r := gin.Default()
	r.Use(middlewares.RateLimiter(), middlewares.CheckJWT(), middlewares.ErrHandler())
	// 构造Service
	memberService := service.NewMemberService(db)
	// 构造Handler
	memberHandler := handler.NewMemberHandler(memberService)

	apiv1Admin := r.Group("/api/v1/")
	memberPlan := apiv1Admin.Group("/member_plan")
	{
		memberPlan.POST("", memberHandler.AddNewMemberPlan)
		memberPlan.PUT("/:plan_id", memberHandler.UpdateMemberPlan)
	}
	memberList := apiv1Admin.Group("/member_list")
	{
		memberList.GET("")
		memberList.POST("")
		memberList.PUT("")
	}
	return r
}

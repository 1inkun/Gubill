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
	r := gin.Default()
	gin.SetMode(ginMode)
	// 中间件
	r.Use(middlewares.RateLimiter(), middlewares.ErrHandler())

	userService := service.NewUserService(db)
	signService := service.NewSignService(db)

	userHandler := handler.NewUserHandler(userService)
	sighHandler := handler.NewSignHandler(signService)

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
	return r
}

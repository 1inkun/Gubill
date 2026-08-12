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
	r.Use(middlewares.ErrHandler())

	userService := service.NewUserServer(db)

	userHandler := handler.NewUserHandler(userService)

	apiv1 := r.Group("/api/v1")
	user := apiv1.Group("/user")
	{
		user.POST("/login", userHandler.Login)
		user.POST("/register", userHandler.Register)
	}
	return r
}

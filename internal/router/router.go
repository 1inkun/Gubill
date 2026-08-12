package router

import (
	"github.com/1inkun/Gubill/internal/config"
	"github.com/1inkun/Gubill/internal/handler"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(c *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()
	gin.SetMode(c.Gin.Mode)

	userService := service.NewUserServer(db)

	userHandler := handler.NewUserHandler(userService)

	apiv1 := r.Group("/api/v1")
	user := apiv1.Group("/user")
	{
		user.POST("/login", userHandler.Login)
		user.POST("/register")
	}
	return r
}

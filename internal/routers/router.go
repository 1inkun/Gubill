package routers

import (
	"database/sql"
	"net/http"

	"github.com/1inkun/billApp/internal/handlers"
	"github.com/1inkun/billApp/internal/middlewares"
	"github.com/1inkun/billApp/internal/repositories"
	"github.com/1inkun/billApp/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	_ "github.com/1inkun/billApp/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title		billApp API
//	@version	1.0

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/api/v1
func NewRouter(db *sql.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()
	// 注册中间件
	r.Use(middlewares.SetRequestId(), middlewares.ErrHandler())
	// 测试接口
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": "Hello World",
		})
	})
	// 初始化仓储层
	userQuery := repositories.NewUserQuery(db, rdb)
	orderQuery := repositories.NewOrderQuery(db, rdb)
	statusQuery := repositories.NewStatusQuery(db, rdb)
	goodQuery := repositories.NewGoodQuery(db, rdb)
	memberQuery := repositories.NewMemberQuery(db, rdb)
	payQuery := repositories.NewPayQuery(db, rdb)

	// 初始化服务层
	userService := services.NewUserService(userQuery)
	orderService := services.NewOrderService(orderQuery)
	statusService := services.NewStatusService(statusQuery)
	goodService := services.NewGoodService(goodQuery)
	memberService := services.NewMemberService(memberQuery)
	payService := services.NewPayService(payQuery)

	// 初始化控制器层
	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)
	statusHandler := handlers.NewStatusService(statusService)
	goodHandler := handlers.NewGoodHandler(goodService)
	memberHandler := handlers.NewMemberHandler(memberService)
	payHandler := handlers.NewPayHandler(payService)

	apiv1 := r.Group("/api/v1")
	// 登录相关的接口
	{
		apiv1.POST("/login", userHandler.Login)
		apiv1.POST("/register", userHandler.Register)
		apiv1.POST("/reflash", userHandler.Reflash)
	}
	// 用户相关接口
	// 签到相关接口
	order := apiv1.Group("/order")
	order.Use(middlewares.Auth(rdb))
	{
		order.POST("", orderHandler.NewOrder)
		order.GET("", orderHandler.GetOrderData)
		order.PUT("/:order_id", orderHandler.FinishOrder)
	}
	status := apiv1.Group("/status")
	{
		status.GET("", statusHandler.GetSginStatus)
	}
	// 会员相关接口
	good := apiv1.Group("/good")
	{
		good.GET("", goodHandler.GetGoods)
		good.POST("/:good_id", middlewares.Auth(rdb), goodHandler.BuyGoods)
	}
	member := apiv1.Group("/member")
	member.Use(middlewares.Auth(rdb))
	{
		member.GET("", memberHandler.GetMemberData)
	}
	// 支付相关的接口
	pay := apiv1.Group("/pay")
	pay.Use(middlewares.Auth(rdb))
	{
		pay.GET("", payHandler.GetPayOrderData)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

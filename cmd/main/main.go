package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/1inkun/Gubill/internal/config"
	"github.com/1inkun/Gubill/internal/payment"
	"github.com/1inkun/Gubill/internal/repository"
	"github.com/1inkun/Gubill/internal/router"
	"github.com/1inkun/Gubill/internal/service"
)

func main() {
	config.InitConfig()
	if os.Getenv("JWTSalt") == "" {
		log.Fatalln("缺少 JWT 密钥，请在 cmd/main/config.yaml 的 JWT.Salt 中配置")
	}
	db, err := repository.InitDatabaseConnect()
	if err != nil {
		log.Fatalf("数据库连接错误:%s", err.Error())
	}
	// 构造支付服务（网关注入点）
	expireMinutes, _ := strconv.ParseInt(os.Getenv("PayExpireMinutes"), 10, 64)
	// TODO(支付接入)：实现 internal/payment.Gateway（微信/支付宝）后，在此创建实例并注入：
	//   gateway := payment.NewWechatGateway(/* 商户配置 */)   // 或支付宝实现
	// 未接入前保持 nil：结算接口将返回"支付网关未配置"，管理端仍可手工确认收款/退款。
	var gateway payment.Gateway
	paymentService := service.NewPaymentService(db, gateway, expireMinutes)
	// 后台定时作废过期支付单
	sweeperCtx, sweeperCancel := context.WithCancel(context.Background())
	defer sweeperCancel()
	go paymentService.StartExpireSweeper(sweeperCtx, 5*time.Minute)

	server := http.Server{
		Addr:    os.Getenv("ServerAddr"),
		Handler: router.InitRouter(db, paymentService),
	}
	serverAdmin := http.Server{
		Addr:    os.Getenv("ServerAdminAddr"),
		Handler: router.InitAdminRouter(db, paymentService),
	}
	// 用户侧服务
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("用户服务出错")
		}
	}()
	go func() {
		if err := serverAdmin.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("管理服务出错")
		}
	}()
	// 优雅关停
	// 等待中断信号，以便在5秒后优雅关停服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	if err := serverAdmin.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

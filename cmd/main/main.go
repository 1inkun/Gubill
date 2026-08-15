package main

import (
	"log"
	"net/http"
	"os"

	"github.com/1inkun/Gubill/internal/config"
	"github.com/1inkun/Gubill/internal/repository"
	"github.com/1inkun/Gubill/internal/router"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	config.InitConfig()
	db, err := repository.InitDatabaseConnect()
	if err != nil {
		log.Fatalf("数据库连接错误:%s", err.Error())
	}
	server := http.Server{
		Addr:    os.Getenv("Addr"),
		Handler: router.InitRouter(db),
	}
	serverAdmin := http.Server{
		Addr:    ":8081",
		Handler: router.InitAdminRouter(db),
	}
	// 用户侧服务
	g.Go(func() error {
		return server.ListenAndServe()
	})
	// 管理侧服务
	g.Go(func() error {
		return serverAdmin.ListenAndServe()
	})
	if err := g.Wait(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	// 优雅关停
	// 等待中断信号，以便在5秒后优雅关停服务器
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	// log.Println("Shutdown Server ...")
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := server.Shutdown(ctx); err != nil {
	// 	log.Println("Server Shutdown:", err)
	// }
	// if err := serverAdmin.Shutdown(ctx); err != nil {
	// 	log.Println("Server Shutdown:", err)
	// }
	// log.Println("Server exiting")
}

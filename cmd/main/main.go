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
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

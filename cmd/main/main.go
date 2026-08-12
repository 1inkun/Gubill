package main

import (
	"log"
	"net/http"

	"github.com/1inkun/Gubill/internal/config"
	"github.com/1inkun/Gubill/internal/repository"
	"github.com/1inkun/Gubill/internal/router"
)

func main() {
	config := config.InitConfig()
	db, err := repository.InitDatabaseConnect(config)
	if err != nil {
		log.Fatalf("数据库连接错误:%s", err.Error())
	}
	_ = db
	r := router.InitRouter(config, db)
	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("服务器出错")
	}
}

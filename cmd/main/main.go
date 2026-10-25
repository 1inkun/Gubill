package main

import (
	"log"
	"net/http"

	"github.com/1inkun/billApp/config"
	"github.com/1inkun/billApp/internal/repositories/mysql"
	"github.com/1inkun/billApp/internal/repositories/redis"
	"github.com/1inkun/billApp/internal/routers"
	"golang.org/x/sync/errgroup"
)

func main() {
	// 读取配置文件
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln("出现致命错误，服务启动失败")
	}
	// log.Println(cfg)
	// 连接数据库
	db, err := mysql.NewDBConection(cfg.MySql)
	if err != nil {
		log.Fatalln("出现致命错误，服务启动失败")
	}
	defer db.Close()
	// 连接Redis
	rdb, err := redis.NewRSConection(cfg.Redis)
	if err != nil {
		log.Fatalln("出现致命错误，服务启动失败")
	}
	defer rdb.Close()
	// 初始化路由
	r := routers.NewRouter(db, rdb)
	serverMux := http.Server{
		Addr:    cfg.Server.Addr,
		Handler: r,
	}
	// 启动服务器
	var eg = errgroup.Group{}
	eg.Go(func() error {
		log.Println("服务启动成功，端口:8080")
		err := serverMux.ListenAndServe()
		if err != nil {
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		log.Fatalln("出现致命错误，服务器已关闭!!")
	}
}

package redis

import (
	"context"
	"log"

	"github.com/1inkun/billApp/config"
	"github.com/redis/go-redis/v9"
)

func NewRSConection(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Host + ":" + cfg.Port,
		DB:   0,
	})
	ctx := context.Background()
	// 测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	// 数据库连接成功
	log.Println("Redis连接成功!")
	return rdb, nil
}

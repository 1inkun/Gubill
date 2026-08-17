package config

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SQLite      SQLite      `yaml:"SQLite"`
	Gin         Gin         `yaml:"Gin"`
	JWT         JWT         `yaml:"JWT"`
	Billing     Billing     `yaml:"Billing"`
	Pay         PayConfig   `yaml:"Pay"`
	Server      Server      `yaml:"Server"`
	ServerAdmin ServerAdmin `yaml:"ServerAdmin"`
}

type SQLite struct {
	Url string `yaml:"Url"`
}

type Gin struct {
	Mode string `yaml:"Mode"`
}

// JWT 签名配置（config.yaml 为唯一配置源，启动时由 config 包写入环境变量）
type JWT struct {
	Salt string `yaml:"Salt"`
}

type Server struct {
	Addr          string `yaml:"Addr"`
	PublicBaseUrl string `yaml:"PublicBaseUrl"`
}

type ServerAdmin struct {
	Addr string `yaml:"Addr"`
}

// Billing 计费相关配置
type Billing struct {
	SinglePrice int64 `yaml:"SinglePrice"` // 签到单价（单位：分）
}

// PayConfig 支付相关配置
type PayConfig struct {
	ExpireMinutes int64 `yaml:"ExpireMinutes"` // 支付单有效期（分钟）
}

func InitConfig() {
	var config Config
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("配置文件读取错误(%s):%s", configPath, err.Error())
	}
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("配置文件解析错误(%s):%s", configPath, err.Error())
	}
	// 遵循上游约定：config.yaml 为唯一配置源，读取后写入环境变量
	os.Setenv("GinMode", config.Gin.Mode)
	os.Setenv("SQLiteUrl", config.SQLite.Url)
	os.Setenv("JWTSalt", config.JWT.Salt)
	os.Setenv("ServerAddr", config.Server.Addr)
	os.Setenv("ServerAdminAddr", config.ServerAdmin.Addr)
	os.Setenv("PublicBaseUrl", config.Server.PublicBaseUrl)
	os.Setenv("SinglePrice", fmtInt(config.Billing.SinglePrice))
	os.Setenv("PayExpireMinutes", fmtInt(config.Pay.ExpireMinutes))
}

func fmtInt(v int64) string {
	if v == 0 {
		return ""
	}
	return strconv.FormatInt(v, 10)
}

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
	// 写入环境变量：已存在的外部环境变量优先，yaml 仅作为默认值
	setDefaultEnv := func(key, val string) {
		if _, ok := os.LookupEnv(key); !ok {
			os.Setenv(key, val)
		}
	}
	setDefaultEnv("GinMode", config.Gin.Mode)
	setDefaultEnv("SQLiteUrl", config.SQLite.Url)
	setDefaultEnv("ServerAddr", config.Server.Addr)
	setDefaultEnv("ServerAdminAddr", config.ServerAdmin.Addr)
	setDefaultEnv("PublicBaseUrl", config.Server.PublicBaseUrl)
	setDefaultEnv("SinglePrice", fmtInt(config.Billing.SinglePrice))
	setDefaultEnv("PayExpireMinutes", fmtInt(config.Pay.ExpireMinutes))
}

func fmtInt(v int64) string {
	if v == 0 {
		return ""
	}
	return strconv.FormatInt(v, 10)
}

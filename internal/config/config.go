package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SQLite SQLite `yaml:"SQLite"`
	Gin    Gin    `yaml:"Gin"`
	JWT    JWT    `yaml:"JWT"`
}

type SQLite struct {
	Url string `yaml:"Url"`
}

type Gin struct {
	Mode string `yaml:"Mode"`
}

type JWT struct {
	Salt string `yaml:"Salt"`
}

func InitConfig() {
	var config Config
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("配置文件读取错误:%s", err.Error())
	}
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("配置文件读取错误:%s", err.Error())
	}
	log.Println(config)
	// 设置环境变量
	os.Setenv("GinMode", config.Gin.Mode)
	os.Setenv("SQLiteUrl", config.SQLite.Url)
	os.Setenv("JWTSalt", config.JWT.Salt)
}

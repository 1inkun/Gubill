package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SQLite SQLite `yaml:"SQLite"`
	Gin    Gin    `yaml:"Gin"`
}

type SQLite struct {
	Url string `yaml:"Url"`
}

type Gin struct {
	Mode string `yaml:"Mode"`
}

func InitConfig() *Config {
	var config Config
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("配置文件读取错误:%s", err.Error())
	}
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("配置文件读取错误:%s", err.Error())
	}
	return &config
}

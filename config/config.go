package config

import (
	"log"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	MySql  MySqlConfig  `yaml:"MySql"`
	Redis  RedisConfig  `yaml:"Redis"`
	Token  TokenConfig  `yaml:"Token"`
	Server ServerConfig `yaml:"Server"`
}

type MySqlConfig struct {
	Host         string `yaml:"Host"`
	Port         string `yaml:"Port"`
	User         string `yaml:"User"`
	Password     string `yaml:"Password"`
	Database     string `yaml:"Database"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
}

type RedisConfig struct {
	Host string `yaml:"Host"`
	Port string `yaml:"Port"`
}

type TokenConfig struct {
	JwtSecret  string `yaml:"JwtSecret"`
	JwtExpTime string `yaml:"JwtExpTime"`
}

type ServerConfig struct {
	Addr  string `yaml:"Addr"`
	Price string `yaml:"Price"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}
	file, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Printf("读取配置文件时出现错误:%s", err.Error())
		return nil, err
	}
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Printf("解析配置文件时出现错误:%s", err.Error())
		return nil, err
	}
	// log.Println(cfg)
	// 设置部分环境变量
	os.Setenv("secret", cfg.Token.JwtSecret)
	os.Setenv("exptime", cfg.Token.JwtExpTime)
	os.Setenv("price", cfg.Server.Price)
	return &cfg, nil
}

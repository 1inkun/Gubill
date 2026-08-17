// gubill-cli 提供运维命令，目前支持 create-admin。
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/1inkun/Gubill/internal/config"
	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/repository"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "create-admin" {
		fmt.Fprintln(os.Stderr, "用法: gubill create-admin --username <用户名> --password <密码> [--email <邮箱>]")
		os.Exit(2)
	}
	fs := flag.NewFlagSet("create-admin", flag.ExitOnError)
	username := fs.String("username", "", "管理员用户名")
	password := fs.String("password", "", "管理员密码")
	email := fs.String("email", "", "管理员邮箱（可选）")
	_ = fs.Parse(os.Args[2:])

	if *username == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "用户名和密码不能为空")
		os.Exit(2)
	}
	if len(*password) < 6 {
		fmt.Fprintln(os.Stderr, "密码长度至少 6 位")
		os.Exit(2)
	}

	config.InitConfig()
	db, err := repository.InitDatabaseConnect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "数据库连接失败:%s\n", err.Error())
		os.Exit(1)
	}

	passwordHash, err := utils.GenNewPasswdHash(*password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "密码加密失败:%s\n", err.Error())
		os.Exit(1)
	}

	emailValue := *email
	if emailValue == "" {
		emailValue = *username + "@local"
	}

	ctx := context.Background()
	err = db.Transaction(func(tx *gorm.DB) error {
		users, e := gorm.G[models.User](tx).Where("username = ?", *username).Limit(1).Find(ctx)
		if e != nil {
			return e
		}
		if len(users) > 0 {
			// 已存在则升级为管理员并重置密码
			e = tx.Model(&models.User{}).Where("uuid = ?", users[0].UUID).Updates(map[string]any{
				"role":          "Admin",
				"password_hash": passwordHash,
				"email":         emailValue,
			}).Error
			return e
		}
		newUser := models.User{
			UserName:     *username,
			NickName:     "管理员",
			Email:        emailValue,
			PasswordHash: passwordHash,
			Role:         "Admin",
			RegisterIP:   "cli",
			RegisterDate: time.Now().Unix(),
		}
		return gorm.G[models.User](tx).Create(ctx, &newUser)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建管理员失败:%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("管理员 %s 创建/更新成功\n", *username)
}

package service

import (
	"context"
	"errors"
	"log"

	"github.com/1inkun/Gubill/internal/middlewares"
	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserServer(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

var (
	ErrNoSuchUser        = errors.New("该用户不存在")
	ErrDatabaseErr       = errors.New("数据库访问出错")
	ErrFailGenJWT        = errors.New("生成JWT令牌时遇到错误")
	ErrUserAlreadyExist  = errors.New("该用户名已被注册")
	ErrFailGenPasswdHash = errors.New("生成密码哈希失败")
)

func (s *UserService) Login(ctx context.Context, username string, password string) (string, error) {
	var userDatas []models.User
	userDatas, err := gorm.G[models.User](s.db).Where("username = ?", username).Limit(1).Find(ctx)
	if err != nil {
		return "", ErrDatabaseErr
	}
	if len(userDatas) == 0 {
		return "", ErrNoSuchUser
	}
	userData := userDatas[0]
	// 比对密码哈希
	if _, err := utils.CheckPassword(userData.PasswordHash, password); err != nil {
		return "", err
	}
	// 登录成功,生成JWT
	tokenString, err := middlewares.GenNewJWT(userData)
	if err != nil {
		return "", ErrFailGenJWT
	}
	return tokenString, nil
}

func (s *UserService) Register(ctx context.Context, username string, nickname string, password string, email string) (string, error) {
	var checkData []models.User
	var newData models.User
	passwordHash, err := utils.GenNewPasswdHash(password)
	if err != nil {
		log.Println(err.Error())
		return "", ErrFailGenPasswdHash
	}
	newData.UserName = username
	newData.NickName = nickname
	newData.PasswordHash = passwordHash
	newData.Email = email
	err = s.db.Transaction(func(tx *gorm.DB) error {
		var transactionErr error
		// 先检查该用户名称数据是否存在
		checkData, transactionErr = gorm.G[models.User](tx).Where("username = ?", username).Find(ctx)
		if err != nil {
			return transactionErr
		}
		if len(checkData) != 0 {
			return ErrUserAlreadyExist
		}
		// 再进行新数据插入
		transactionErr = gorm.G[models.User](tx).Create(ctx, &newData)
		if transactionErr != nil {
			return transactionErr
		}
		return nil
	})
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return newData.UUID, nil
}

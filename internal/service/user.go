package service

import (
	"context"
	"log"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/repository"
	"github.com/1inkun/Gubill/internal/utils"
)

type UserService struct {
	userQuery *repository.UserQuery
}

func NewUserService(userQuery *repository.UserQuery) *UserService {
	return &UserService{userQuery: userQuery}
}

func (s *UserService) Login(ctx context.Context, bodyData models.LoginData) (string, error) {
	var userData models.User
	userData, err := s.userQuery.Login(ctx, bodyData.Username, bodyData.Password)
	if err != nil {
		return "", err
	}
	tokenString, err := utils.GenNewJWT(userData)
	if err != nil {
		ErrFailGenJWT := models.NewInternalError(500, "登录失败", err)
		return "", ErrFailGenJWT
	}
	return tokenString, nil
}

func (s *UserService) Register(ctx context.Context, data models.RegisterData) (string, error) {
	// var newData models.User
	passwordHash, err := utils.GenNewPasswdHash(data.Password)
	if err != nil {
		log.Println(err.Error())
		ErrFailGenPasswdHash := models.NewInternalError(500, "注册失败", err)
		return "", ErrFailGenPasswdHash
	}
	// newData.UserName = username
	// newData.NickName = nickname
	// newData.PasswordHash = passwordHash
	// newData.Email = email
	// 直接插入数据,利用唯一性标签处理用户名和邮箱占用
	// err = gorm.G[models.User](s.db).Create(ctx, &newData)
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrDuplicatedKey) {
	// 		return "", ErrWrongEmaiOrUserName
	// 	}
	// 	return "", models.NewDatabaseErr(err)
	// }
	userId, err := s.userQuery.Register(ctx, data.UserName, data.NickName, passwordHash, data.Email)
	if err != nil {
		return "", err
	}
	return userId, nil
}

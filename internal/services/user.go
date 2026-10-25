package services

import (
	"context"
	"log"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/repositories"
	"github.com/1inkun/billApp/internal/utils"
)

type UserService struct {
	userQuery *repositories.UserQuery
}

func NewUserService(userQuery *repositories.UserQuery) *UserService {
	return &UserService{
		userQuery: userQuery,
	}
}

func (s *UserService) Register(ctx context.Context, data models.RegisterRequest) error {
	passwdHash, err := utils.GeneratePasswdHash(data.Password)
	if err != nil {
		log.Println(err.Error())
		// 属于内部错误
		return err
	}
	registerData := models.RegisterData{
		Username:   data.Username,
		Nickname:   data.Nickname,
		PasswdHash: passwdHash,
		Email:      data.Email,
	}
	err = s.userQuery.Register(ctx, registerData)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(ctx context.Context, logInfo models.LoginRequest) (string, error) {
	passwdHash, err := s.userQuery.GetUserPasswdHashByUername(ctx, logInfo.Username)
	if err != nil {
		return "", err
	}
	// 验证密码
	if ok := utils.ComparePasswdHash(passwdHash, logInfo.Password); !ok {
		return "", models.ErrWrongPasswd
	}
	// 获取JWTClaim
	claims, err := s.userQuery.GetClaimsData(ctx, logInfo.Username)
	if err != nil {
		return "", err
	}
	// 生成JWT
	tokenStr, err := utils.GenerateJWT(claims)
	if err != nil {
		return "", err
	}
	// 将token存入redis中
	err = s.userQuery.AddTokenToRedis(ctx, logInfo.Username, tokenStr)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

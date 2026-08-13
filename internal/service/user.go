package service

import (
	"context"
	"log"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

var (
	ErrNoSuchUser          = models.NewBusinessError(404, "用户名或密码错误")
	ErrDatabaseErr         = models.NewInternalError(500, "数据库访问出错")
	ErrFailGenJWT          = models.NewInternalError(500, "生成令牌时遇到错误")
	ErrUserAlreadyExist    = models.NewBusinessError(400, "该用户名已被注册")
	ErrFailGenPasswdHash   = models.NewInternalError(500, "生成密码哈希错误")
	ErrWrongPasswd         = models.NewBusinessError(401, "密码错误")
	ErrWrongEmaiOrUserName = models.NewBusinessError(401, "用户名或邮箱已被占用")
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
		return "", ErrWrongPasswd
	}
	// 登录成功,生成JWT
	tokenString, err := utils.GenNewJWT(userData)
	if err != nil {
		return "", ErrFailGenJWT
	}
	return tokenString, nil
}

func (s *UserService) Register(ctx context.Context, username string, nickname string, password string, email string) (string, error) {
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
	// 直接插入数据,利用唯一性标签处理用户名和邮箱占用
	err = gorm.G[models.User](s.db).Create(ctx, &newData)
	if err != nil {
		return "", ErrWrongEmaiOrUserName
	}
	return newData.UUID, nil
}

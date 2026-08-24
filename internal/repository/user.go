package repository

import (
	"context"
	"log"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

type UserQuery struct {
	db *gorm.DB
}

func NewUserQuery(db *gorm.DB) *UserQuery {
	return &UserQuery{db: db}
}

var (
	ErrNoSuchUser          = models.NewBusinessError(401, "用户名或密码错误")
	ErrUserAlreadyExist    = models.NewBusinessError(400, "该用户名已被注册")
	ErrWrongPasswd         = models.NewBusinessError(401, "用户名或密码错误")
	ErrWrongEmaiOrUserName = models.NewBusinessError(401, "用户名或邮箱已被占用")
)

func (q *UserQuery) Login(ctx context.Context, username string, password string) (models.User, error) {
	result, err := gorm.G[models.User](q.db).Where("username = ?", username).First(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, ErrNoSuchUser
		}
		log.Println(err)
		return models.User{}, models.NewDatabaseErr(err)
	}
	// 对比密码哈希
	if _, err := utils.CheckPassword(result.PasswordHash, password); err != nil {
		return models.User{}, ErrWrongPasswd
	}
	return result, nil
}

func (q *UserQuery) Register(ctx context.Context, username string, nickname string, passwdHash string, email string) (string, error) {
	var newData = models.User{
		UserName:     username,
		NickName:     nickname,
		PasswordHash: passwdHash,
		Email:        email,
	}
	// 直接插入数据,利用唯一性标签处理用户名和邮箱占用
	err := gorm.G[models.User](q.db).Create(ctx, &newData)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return "", ErrWrongEmaiOrUserName
		}
		return "", models.NewDatabaseErr(err)
	}
	return newData.UUID, nil
}

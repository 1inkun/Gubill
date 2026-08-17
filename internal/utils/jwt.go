package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserId   string
	Username string
	Nickname string
	Role     string
	jwt.RegisteredClaims
}

func GenNewJWT(userData models.User) (string, error) {
	signingKey := []byte(os.Getenv("JWT_SALT"))
	claims := JWTClaims{
		userData.UUID,
		userData.UserName,
		userData.NickName,
		userData.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
			Issuer:    "Server",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(signingKey)
	if err != nil {
		log.Printf("生成JWT时出错:%s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(JWTString string) (JWTClaims, error) {
	token, err := jwt.ParseWithClaims(JWTString, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("错误的JWT算法")
		}
		return []byte(os.Getenv("JWT_SALT")), nil
	})
	if err != nil {
		return JWTClaims{}, err
	} else if claims, ok := token.Claims.(*JWTClaims); ok {
		return *claims, nil
	} else {
		return JWTClaims{}, errors.New("处理JWT载荷错误")
	}
}

package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/1inkun/billApp/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	// Exp      int64
	// Iat      int64
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(claim models.JwtClaims) (string, error) {
	secret := os.Getenv("secret")
	exptime, err := strconv.Atoi(os.Getenv("exptime"))
	if err != nil {
		log.Panicf("读取配置文件时出现错误:%s", err)
		return "", err
	}
	claims := Claims{
		claim.UserId,
		claim.Username,
		claim.Role,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add((time.Duration(exptime) * time.Hour))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenStr, nil
}

func ParseJWT(tokenStr string) (models.JwtClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		hmacSampleSecret := []byte("你好")
		return hmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Println(err)
		return models.JwtClaims{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var data models.JwtClaims
		// log.Println(claims["user_id"])
		data.Role = claims["role"].(string)
		data.Username = claims["username"].(string)
		userId := claims["user_id"].(float64)
		data.UserId = int64(userId)
		return data, nil
	}
	return models.JwtClaims{}, nil
}

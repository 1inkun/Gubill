package middlewares

import (
	"log"
	"os"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Username string
	Nickname string
	Role     string
	jwt.RegisteredClaims
}

func GenNewJWT(userData models.User) (string, error) {
	signingKey := []byte(os.Getenv("JWTSalt"))
	claims := JWTClaims{
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

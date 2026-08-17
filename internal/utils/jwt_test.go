package utils

import (
	"os"
	"testing"
	"time"

	"github.com/1inkun/Gubill/internal/models"
)

func TestJWTRoundTrip(t *testing.T) {
	os.Setenv("JWTSalt", "test-jwt-salt")
	user := models.User{UserName: "tester", NickName: "T", Role: "Admin"}
	user.UUID = "u-1"
	token, err := GenNewJWT(user)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %s", err.Error())
	}
	claims, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("解析 JWT 失败: %s", err.Error())
	}
	if claims.UserId != "u-1" || claims.Role != "Admin" || claims.Username != "tester" {
		t.Errorf("claims 不符: %+v", claims)
	}
	if !claims.ExpiresAt.After(time.Now()) {
		t.Error("JWT 有效期应大于当前时间")
	}
}

func TestJWTParsingWithWrongSalt(t *testing.T) {
	os.Setenv("JWTSalt", "salt-a")
	user := models.User{UserName: "u", Role: "User"}
	user.UUID = "u-2"
	token, err := GenNewJWT(user)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %s", err.Error())
	}
	os.Setenv("JWTSalt", "salt-b")
	if _, err := ParseJWT(token); err == nil {
		t.Error("使用错误密钥应解析失败")
	}
}

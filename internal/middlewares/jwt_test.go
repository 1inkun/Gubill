package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestCheckJWT(t *testing.T) {
	os.Setenv("JWT_SALT", "test-salt")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CheckJWT())
	r.GET("/x", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"userId": c.GetString("userId"), "role": c.GetString("Role")})
	})

	// 无 Authorization 头
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("无 token 应 403, got %d", w.Code)
	}

	user := models.User{UserName: "u", Role: "Admin"}
	user.UUID = "u1"
	token, err := utils.GenNewJWT(user)
	if err != nil {
		t.Fatal(err)
	}

	// Bearer 前缀
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"userId":"u1"`) {
		t.Errorf("Bearer 前缀解析失败: %d %s", w.Code, w.Body.String())
	}

	// 裸 token
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("裸 token 解析失败: %d %s", w.Code, w.Body.String())
	}

	// 过期 token
	expired, err := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.JWTClaims{
		UserId: "u1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}).SignedString([]byte(os.Getenv("JWT_SALT")))
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+expired)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("过期 token 应 403, got %d", w.Code)
	}
}

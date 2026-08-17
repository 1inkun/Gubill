package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

func TestAdminLoginPage(t *testing.T) {
	_, r := newAdminRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "管理员登录") {
		t.Fatalf("登录页异常: %d %s", w.Code, w.Body.String())
	}
}

func TestAdminLoginWrongPassword(t *testing.T) {
	db, r := newAdminRouter(t)
	ctx := context.Background()
	hash, _ := utils.GenNewPasswdHash("admin1234")
	admin := models.User{UserName: "boss", NickName: "管理员", Email: "boss@test.com", PasswordHash: hash, Role: "Admin"}
	if err := gorm.G[models.User](db).Create(ctx, &admin); err != nil {
		t.Fatal(err)
	}
	form := url.Values{"username": {"boss"}, "password": {"wrongpass"}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "用户名或密码错误") {
		t.Fatalf("错误密码应留在登录页并提示: %d %s", w.Code, w.Body.String())
	}
}

func TestAdminLoginSetsCookieAndAuth(t *testing.T) {
	db, r := newAdminRouter(t)
	ctx := context.Background()
	hash, _ := utils.GenNewPasswdHash("admin1234")
	admin := models.User{UserName: "boss", NickName: "管理员", Email: "boss@test.com", PasswordHash: hash, Role: "Admin"}
	if err := gorm.G[models.User](db).Create(ctx, &admin); err != nil {
		t.Fatal(err)
	}

	// 未登录访问管理页 → 302 到登录页
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusFound || w.Header().Get("Location") != "/admin/login" {
		t.Fatalf("未登录应重定向: %d %s", w.Code, w.Header().Get("Location"))
	}

	// 登录成功 → 302 且种 Cookie
	form := url.Values{"username": {"boss"}, "password": {"admin1234"}}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Fatalf("登录应重定向: %d %s", w.Code, w.Body.String())
	}
	cookies := w.Result().Cookies()
	var token string
	for _, ck := range cookies {
		if ck.Name == "gubill_token" {
			token = ck.Value
		}
	}
	if token == "" {
		t.Fatal("登录后应种 gubill_token Cookie")
	}

	// 带 Cookie 访问管理页 → 200 概览
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "gubill_token", Value: token})
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "概览") {
		t.Fatalf("带 Cookie 访问概览失败: %d %s", w.Code, w.Body.String())
	}

	// 支付记录页可打开
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/admin/pays", nil)
	req.AddCookie(&http.Cookie{Name: "gubill_token", Value: token})
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "支付记录") {
		t.Fatalf("支付记录页异常: %d %s", w.Code, w.Body.String())
	}
}

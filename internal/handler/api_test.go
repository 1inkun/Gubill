package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/router"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/1inkun/Gubill/internal/testutil"
	"github.com/1inkun/Gubill/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SALT", "test-salt")
	os.Setenv("SinglePrice", "500")
	os.Setenv("PayExpireMinutes", "30")
	os.Setenv("GinMode", "test")
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newUserRouter(t *testing.T) (*gorm.DB, http.Handler) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	ps := service.NewPaymentService(db, testutil.FakeGatewayInstance(), 30)
	return db, router.InitRouter(db, ps)
}

func newAdminRouter(t *testing.T) (*gorm.DB, http.Handler) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	ps := service.NewPaymentService(db, testutil.FakeGatewayInstance(), 30)
	return db, router.InitAdminRouter(db, ps)
}

func doJSON(t *testing.T, r http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func apiResult(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("解析响应失败(%s): %s", w.Body.String(), err.Error())
	}
	return body
}

func registerAndLogin(t *testing.T, r http.Handler, username, password string) string {
	t.Helper()
	w := doJSON(t, r, "POST", "/api/v1/user/register",
		`{"username":"`+username+`","nickname":"测试","password":"`+password+`","email":"`+username+`@test.com"}`,
		"")
	if w.Code != http.StatusOK {
		t.Fatalf("注册失败: %d %s", w.Code, w.Body.String())
	}
	w = doJSON(t, r, "POST", "/api/v1/user/login",
		`{"username":"`+username+`","password":"`+password+`"}`, "")
	if w.Code != http.StatusOK {
		t.Fatalf("登录失败: %d %s", w.Code, w.Body.String())
	}
	data := apiResult(t, w)
	token, _ := data["data"].(map[string]any)["token"].(string)
	if token == "" {
		t.Fatalf("未返回 token: %s", w.Body.String())
	}
	return token
}

func TestUserPayFlow(t *testing.T) {
	db, r := newUserRouter(t)
	_ = db
	token := registerAndLogin(t, r, "user1", "pass1234")
	bearer := "Bearer " + token

	// 开始签到
	w := doJSON(t, r, "POST", "/api/v1/sign", "", bearer)
	if w.Code != http.StatusOK {
		t.Fatalf("开始签到失败: %d %s", w.Code, w.Body.String())
	}
	signId := apiResult(t, w)["data"].(map[string]any)["signId"].(string)

	// 结算生成支付单
	w = doJSON(t, r, "PUT", "/api/v1/sign/"+signId, "", bearer)
	if w.Code != http.StatusOK {
		t.Fatalf("结算失败: %d %s", w.Code, w.Body.String())
	}
	data := apiResult(t, w)["data"].(map[string]any)
	payId := data["payId"].(string)
	payUrl := data["payUrl"].(string)
	if payId == "" || payUrl == "" {
		t.Fatalf("结算结果异常: %+v", data)
	}

	// 支付单列表可见
	w = doJSON(t, r, "GET", "/api/v1/pay?page=1&page_size=10", "", bearer)
	if w.Code != http.StatusOK {
		t.Fatalf("查询支付记录失败: %d %s", w.Code, w.Body.String())
	}
	results := apiResult(t, w)["data"].(map[string]any)["results"].([]any)
	if len(results) != 1 {
		t.Fatalf("应返回 1 条支付记录, got %d", len(results))
	}
	if int64(results[0].(map[string]any)["status"].(float64)) != models.PayStatusPending {
		t.Errorf("未支付时状态应为待支付")
	}

	// 其他用户无权查看
	token2 := registerAndLogin(t, r, "user2", "pass1234")
	w = doJSON(t, r, "GET", "/api/v1/pay/"+payId, "", "Bearer "+token2)
	if w.Code != http.StatusBadRequest {
		t.Errorf("越权查看应返回 400, got %d", w.Code)
	}
	// 本人可查看
	w = doJSON(t, r, "GET", "/api/v1/pay/"+payId, "", bearer)
	if w.Code != http.StatusOK {
		t.Errorf("本人查看支付单失败: %d", w.Code)
	}
}

func TestAdminPayFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := service.NewPaymentService(db, testutil.FakeGatewayInstance(), 30)
	userR := router.InitRouter(db, ps)
	adminR := router.InitAdminRouter(db, ps)
	ctx := context.Background()

	// 直接创建管理员
	hash, err := utils.GenNewPasswdHash("admin1234")
	if err != nil {
		t.Fatal(err)
	}
	admin := models.User{UserName: "admin", NickName: "管理员", Email: "admin@test.com", PasswordHash: hash, Role: "Admin"}
	if err := gorm.G[models.User](db).Create(ctx, &admin); err != nil {
		t.Fatal(err)
	}
	adminToken, err := utils.GenNewJWT(admin)
	if err != nil {
		t.Fatal(err)
	}
	adminBearer := "Bearer " + adminToken

	// 直接创建普通用户
	customerHash, _ := utils.GenNewPasswdHash("pass1234")
	customer := models.User{UserName: "customer", NickName: "客户", Email: "customer@test.com", PasswordHash: customerHash}
	if err := gorm.G[models.User](db).Create(ctx, &customer); err != nil {
		t.Fatal(err)
	}
	customerToken, err := utils.GenNewJWT(customer)
	if err != nil {
		t.Fatal(err)
	}
	userBearer := "Bearer " + customerToken

	// 新增会员计划
	w := doJSON(t, adminR, "POST", "/api/v1/member_plan",
		`{"name":"月卡","type":"month","value":3000,"des":"一个月"}`, adminBearer)
	if w.Code != http.StatusOK {
		t.Fatalf("新增计划失败: %d %s", w.Code, w.Body.String())
	}
	planId := apiResult(t, w)["data"].(map[string]any)["planId"].(string)

	// 普通用户生成会员订单（用户端服务）
	w = doJSON(t, userR, "POST", "/api/v1/member_plan/"+planId, "", userBearer)
	if w.Code != http.StatusOK {
		t.Fatalf("生成会员订单失败: %d %s", w.Code, w.Body.String())
	}
	orderId := apiResult(t, w)["data"].(map[string]any)["payId"].(string)

	// 结算订单生成支付单
	w = doJSON(t, userR, "POST", "/api/v1/member_order/"+orderId, "", userBearer)
	if w.Code != http.StatusOK {
		t.Fatalf("结算订单失败: %d %s", w.Code, w.Body.String())
	}
	payId := apiResult(t, w)["data"].(map[string]any)["payId"].(string)

	// 普通用户不可调用管理端确认接口
	w = doJSON(t, adminR, "POST", "/api/v1/pay/"+payId+"/confirm", "", userBearer)
	if w.Code != http.StatusForbidden {
		t.Errorf("普通用户确认收款应返回 403, got %d", w.Code)
	}

	// 管理员确认收款
	w = doJSON(t, adminR, "POST", "/api/v1/pay/"+payId+"/confirm", "", adminBearer)
	if w.Code != http.StatusOK {
		t.Fatalf("管理员确认收款失败: %d %s", w.Code, w.Body.String())
	}

	// 管理员查看全部支付记录
	w = doJSON(t, adminR, "GET", "/api/v1/pay?page=1", "", adminBearer)
	if w.Code != http.StatusOK {
		t.Fatalf("查询全部支付记录失败: %d %s", w.Code, w.Body.String())
	}
	results := apiResult(t, w)["data"].(map[string]any)["results"].([]any)
	if len(results) != 1 {
		t.Fatalf("应返回 1 条支付记录, got %d", len(results))
	}

	// 管理员退款
	w = doJSON(t, adminR, "POST", "/api/v1/pay/"+payId+"/refund", "", adminBearer)
	if w.Code != http.StatusOK {
		t.Fatalf("退款失败: %d %s", w.Code, w.Body.String())
	}
	w = doJSON(t, adminR, "GET", "/api/v1/pay/"+payId, "", adminBearer)
	payData := apiResult(t, w)["data"].(map[string]any)
	if int64(payData["status"].(float64)) != models.PayStatusRefunded {
		t.Errorf("退款后状态应为已退款, got %v", payData["status"])
	}
}

func TestRegisterMessageAndBearer(t *testing.T) {
	_, r := newUserRouter(t)
	w := doJSON(t, r, "POST", "/api/v1/user/register",
		`{"username":"newbie","nickname":"新","password":"pass1234","email":"newbie@test.com"}`, "")
	if w.Code != http.StatusOK {
		t.Fatalf("注册失败: %d %s", w.Code, w.Body.String())
	}
	body := apiResult(t, w)
	if body["msg"] != "注册成功" {
		t.Errorf("注册文案应修正为注册成功, got %v", body["msg"])
	}

	// 裸 token 与 Bearer 均可用（登录接口本身无需鉴权，这里验证登录返回 token 格式）
	w = doJSON(t, r, "POST", "/api/v1/user/login", `{"username":"newbie","password":"pass1234"}`, "")
	if w.Code != http.StatusOK {
		t.Fatalf("登录失败: %d %s", w.Code, w.Body.String())
	}
	token := apiResult(t, w)["data"].(map[string]any)["token"].(string)
	if !strings.HasPrefix(token, "eyJ") {
		t.Errorf("token 格式异常: %s", token)
	}
}

// Package testutil 提供测试用的内存数据库与模拟网关。
package testutil

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/payment"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB 创建共享内存 SQLite，并完成全部表迁移。
// 注意：必须限制单连接，否则多个连接会各自持有独立的 :memory: 数据库。
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// 每个测试使用独立的命名内存库，避免测试间数据串扰
	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", name)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("打开内存数据库失败: %s", err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("获取底层连接失败: %s", err.Error())
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(
		&models.User{}, &models.Sign{}, &models.Pay{},
		&models.MemberList{}, &models.MemberOrders{}, &models.MemberPlan{},
	); err != nil {
		t.Fatalf("自动迁移失败: %s", err.Error())
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}

// FakeGateway 仅用于测试的网关实现：
// 返回固定格式的支付地址与流水号，不代表任何真实支付行为。
type FakeGateway struct{}

// Name 实现 payment.Gateway。
func (g *FakeGateway) Name() string { return "fake" }

// CreateOrder 实现 payment.Gateway。
func (g *FakeGateway) CreateOrder(_ context.Context, order *payment.Order) (*payment.CreateResult, error) {
	return &payment.CreateResult{
		PayUrl:         "http://pay.local/" + order.PayId,
		ChannelOrderNo: "fake-order-" + order.PayId,
	}, nil
}

// VerifyNotifySignature 实现 payment.Gateway（测试中不使用回调验签）。
func (g *FakeGateway) VerifyNotifySignature(_ []byte, _ string) bool { return false }

// Refund 实现 payment.Gateway。
func (g *FakeGateway) Refund(_ context.Context, _ *payment.RefundRequest) error { return nil }

// NewTransactionId 实现 payment.Gateway。
func (g *FakeGateway) NewTransactionId() string { return "fake-txn" }

// FakeGatewayInstance 返回测试用网关实例。
func FakeGatewayInstance() *FakeGateway { return &FakeGateway{} }

package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/testutil"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Setenv("JWTSalt", "test-salt")
	os.Setenv("SinglePrice", "500")
	os.Setenv("PayExpireMinutes", "30")
	os.Exit(m.Run())
}

func newTestPaymentService(db *gorm.DB) *PaymentService {
	return NewPaymentService(db, testutil.FakeGatewayInstance())
}

// mustCreatePay 直接写入一条待支付支付单（业务服务结算时也会这样生成 pays）。
func mustCreatePay(t *testing.T, db *gorm.DB, businessType, businessId, userId string, value int64) *models.Pay {
	t.Helper()
	pay := models.Pay{
		UserId:       userId,
		BusinessType: businessType,
		BusinessId:   businessId,
		Value:        value,
		Status:       models.PayStatusPending,
		ExpireTime:   time.Now().Unix() + 30*60,
	}
	if err := gorm.G[models.Pay](db).Create(context.Background(), &pay); err != nil {
		t.Fatalf("创建支付单失败: %s", err.Error())
	}
	return &pay
}

func TestConfirmPaidAndBusinessLink(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db)
	ctx := context.Background()

	sign := models.Sign{
		UserId:  "u1",
		Status:  models.SignStatusPendingPay,
		StartAt: time.Now().Unix() - 3600,
		EndAt:   time.Now().Unix(),
		Value:   100,
	}
	if err := gorm.G[models.Sign](db).Create(ctx, &sign); err != nil {
		t.Fatalf("创建签到单失败: %s", err.Error())
	}
	pay := mustCreatePay(t, db, "sign", sign.UUID, "u1", 100)

	if err := ps.ConfirmPaid(ctx, pay.UUID); err != nil {
		t.Fatalf("确认支付失败: %s", err.Error())
	}
	var payDb models.Pay
	if err := db.First(&payDb, "uuid = ?", pay.UUID).Error; err != nil {
		t.Fatalf("查询支付单失败: %s", err.Error())
	}
	if payDb.Status != models.PayStatusPaid || payDb.PayAt == 0 || payDb.TransactionId == "" {
		t.Errorf("支付单状态异常: %+v", payDb)
	}
	var signDb models.Sign
	if err := db.First(&signDb, "uuid = ?", sign.UUID).Error; err != nil {
		t.Fatalf("查询签到单失败: %s", err.Error())
	}
	if signDb.Status != models.SignStatusFinished {
		t.Errorf("签到单应联动为已完成, got %d", signDb.Status)
	}

	if err := ps.ConfirmPaid(ctx, pay.UUID); err != ErrPayAlreadyPaid {
		t.Errorf("重复确认应报 ErrPayAlreadyPaid, got %v", err)
	}
}

func TestConfirmPaidStateMachine(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db)
	ctx := context.Background()

	// 已过期
	expired := mustCreatePay(t, db, "sign", "b2", "u1", 100)
	if err := db.Model(&models.Pay{}).Where("uuid = ?", expired.UUID).Update("expire_time", time.Now().Unix()-10).Error; err != nil {
		t.Fatal(err)
	}
	if err := ps.ConfirmPaid(ctx, expired.UUID); err != ErrPayExpired {
		t.Errorf("过期单应报 ErrPayExpired, got %v", err)
	}
	var expiredDb models.Pay
	db.First(&expiredDb, "uuid = ?", expired.UUID)
	if expiredDb.Status != models.PayStatusVoid {
		t.Errorf("过期单应作废, got %d", expiredDb.Status)
	}

	// 已作废
	vo := mustCreatePay(t, db, "sign", "b3", "u1", 100)
	if err := db.Model(&models.Pay{}).Where("uuid = ?", vo.UUID).Update("status", models.PayStatusVoid).Error; err != nil {
		t.Fatal(err)
	}
	if err := ps.ConfirmPaid(ctx, vo.UUID); err != ErrPayVoided {
		t.Errorf("作废单应报 ErrPayVoided, got %v", err)
	}

	// 已退款
	rf := mustCreatePay(t, db, "sign", "b4", "u1", 100)
	if err := db.Model(&models.Pay{}).Where("uuid = ?", rf.UUID).Update("status", models.PayStatusRefunded).Error; err != nil {
		t.Fatal(err)
	}
	if err := ps.ConfirmPaid(ctx, rf.UUID); err != ErrPayRefunded {
		t.Errorf("退款单应报 ErrPayRefunded, got %v", err)
	}
}

func TestRefundRules(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db)
	ctx := context.Background()

	pending := mustCreatePay(t, db, "sign", "b9", "u1", 100)
	if err := ps.RefundPay(ctx, pending.UUID); err != ErrPayNotPaid {
		t.Errorf("待支付单退款应报 ErrPayNotPaid, got %v", err)
	}

	paid := mustCreatePay(t, db, "sign", "b10", "u1", 100)
	if err := ps.ConfirmPaid(ctx, paid.UUID); err != nil {
		t.Fatal(err)
	}
	if err := ps.RefundPay(ctx, paid.UUID); err != nil {
		t.Fatalf("退款失败: %s", err.Error())
	}
	var paidDb models.Pay
	db.First(&paidDb, "uuid = ?", paid.UUID)
	if paidDb.Status != models.PayStatusRefunded || paidDb.RefundValue != 100 || paidDb.RefundAt == 0 {
		t.Errorf("退款记录异常: %+v", paidDb)
	}
	if err := ps.RefundPay(ctx, paid.UUID); err != ErrPayRefunded {
		t.Errorf("重复退款应报 ErrPayRefunded, got %v", err)
	}

	vo := mustCreatePay(t, db, "sign", "b11", "u1", 100)
	db.Model(&models.Pay{}).Where("uuid = ?", vo.UUID).Update("status", models.PayStatusVoid)
	if err := ps.RefundPay(ctx, vo.UUID); err != ErrPayVoided {
		t.Errorf("作废单退款应报 ErrPayVoided, got %v", err)
	}
}

func TestCancelMemberOrder(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ms := NewMemberService(db)
	ctx := context.Background()

	plan := models.MemberPlan{Name: "月卡", Type: "month", Value: 3000, Description: "一个月"}
	if err := gorm.G[models.MemberPlan](db).Create(ctx, &plan); err != nil {
		t.Fatal(err)
	}
	order := models.MemberOrders{PlanId: plan.UUID, UserId: "u1", Status: models.MemberOrderStatusPendingPay}
	if err := gorm.G[models.MemberOrders](db).Create(ctx, &order); err != nil {
		t.Fatal(err)
	}
	pay := mustCreatePay(t, db, "member", order.UUID, "u1", 3000)

	n, err := ms.CancelMemberOrder(ctx, "u1", order.UUID)
	if err != nil || n == 0 {
		t.Fatalf("取消订单失败: %v, rows=%d", err, n)
	}
	var orderDb models.MemberOrders
	db.First(&orderDb, "uuid = ?", order.UUID)
	if orderDb.Status != models.MemberOrderStatusCanceled {
		t.Errorf("订单应取消, got %d", orderDb.Status)
	}
	// 取消订单不操作 pays（支付单状态由支付模块独立管理）
	var payDb models.Pay
	db.First(&payDb, "uuid = ?", pay.UUID)
	if payDb.Status != models.PayStatusPending {
		t.Errorf("取消订单不应改动支付单状态, got %d", payDb.Status)
	}

	// 已支付订单不可取消
	order2 := models.MemberOrders{PlanId: plan.UUID, UserId: "u1", Status: models.MemberOrderStatusPaid}
	if err := gorm.G[models.MemberOrders](db).Create(ctx, &order2); err != nil {
		t.Fatal(err)
	}
	if _, err := ms.CancelMemberOrder(ctx, "u1", order2.UUID); err != ErrPaidOrderCannotCancel {
		t.Errorf("已支付订单应报 ErrPaidOrderCannotCancel, got %v", err)
	}
}

func TestSignFinishFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ss := NewSignService(db)
	ps := newTestPaymentService(db)
	ctx := context.Background()

	sign := models.Sign{UserId: "u1", Status: models.SignStatusActive, StartAt: time.Now().Unix() - 3600}
	if err := gorm.G[models.Sign](db).Create(ctx, &sign); err != nil {
		t.Fatal(err)
	}
	payId, err := ss.FinishSignData(ctx, "u1", sign.UUID)
	if err != nil {
		t.Fatalf("结算失败: %s", err.Error())
	}
	if payId == "" {
		t.Fatal("结算应返回支付单 ID")
	}
	var signDb models.Sign
	db.First(&signDb, "uuid = ?", sign.UUID)
	if signDb.Status != models.SignStatusPendingPay {
		t.Errorf("签到单应为待支付, got %d", signDb.Status)
	}
	var payDb models.Pay
	if err := db.Where("uuid = ?", payId).First(&payDb).Error; err != nil {
		t.Fatalf("支付单应存在: %s", err.Error())
	}
	if payDb.Value != 1000 || payDb.BusinessType != models.PayBusinessSign {
		t.Errorf("支付单数据异常: %+v", payDb)
	}

	// 已结算的签到单不可重复结算（上游语义）
	if _, err := ss.FinishSignData(ctx, "u1", sign.UUID); err != ErrSignDataNoExist {
		t.Errorf("重复结算应报 ErrSignDataNoExist, got %v", err)
	}

	// 支付确认后签到联动完成
	if err := ps.ConfirmPaid(ctx, payId); err != nil {
		t.Fatalf("确认支付失败: %s", err.Error())
	}
	db.First(&signDb, "uuid = ?", sign.UUID)
	if signDb.Status != models.SignStatusFinished {
		t.Errorf("签到单应联动为已完成, got %d", signDb.Status)
	}
}

func TestMemberFinishFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ms := NewMemberService(db)
	ps := newTestPaymentService(db)
	ctx := context.Background()

	plan := models.MemberPlan{Name: "周卡", Type: "week", Value: 2000, Description: "一周"}
	if err := gorm.G[models.MemberPlan](db).Create(ctx, &plan); err != nil {
		t.Fatal(err)
	}
	order := models.MemberOrders{PlanId: plan.UUID, UserId: "u1", Status: models.MemberOrderStatusCreated}
	if err := gorm.G[models.MemberOrders](db).Create(ctx, &order); err != nil {
		t.Fatal(err)
	}
	payId, err := ms.FinishMemberOrder(ctx, "u1", order.UUID)
	if err != nil {
		t.Fatalf("会员结算失败: %s", err.Error())
	}
	if payId == "" {
		t.Fatal("结算应返回支付单 ID")
	}
	var orderDb models.MemberOrders
	db.First(&orderDb, "uuid = ?", order.UUID)
	if orderDb.Status != models.MemberOrderStatusPendingPay {
		t.Errorf("订单应为待支付, got %d", orderDb.Status)
	}
	var payDb models.Pay
	db.First(&payDb, "uuid = ?", payId)
	if payDb.Value != 2000 || payDb.BusinessType != models.PayBusinessMember {
		t.Errorf("支付单数据异常: %+v", payDb)
	}
	// 确认支付后订单联动为已支付
	if err := ps.ConfirmPaid(ctx, payId); err != nil {
		t.Fatalf("确认支付失败: %s", err.Error())
	}
	db.First(&orderDb, "uuid = ?", order.UUID)
	if orderDb.Status != models.MemberOrderStatusPaid {
		t.Errorf("订单应联动为已支付, got %d", orderDb.Status)
	}
}

func TestExpireSweeper(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db)
	ctx := context.Background()

	pay := mustCreatePay(t, db, "sign", "b12", "u1", 100)
	if err := db.Model(&models.Pay{}).Where("uuid = ?", pay.UUID).Update("expire_time", time.Now().Unix()-60).Error; err != nil {
		t.Fatal(err)
	}
	n, err := ps.expirePending(ctx)
	if err != nil {
		t.Fatalf("清理失败: %s", err.Error())
	}
	if n != 1 {
		t.Errorf("应作废 1 笔, got %d", n)
	}
	var payDb models.Pay
	db.First(&payDb, "uuid = ?", pay.UUID)
	if payDb.Status != models.PayStatusVoid {
		t.Errorf("过期单应作废, got %d", payDb.Status)
	}
}

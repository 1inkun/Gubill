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
	os.Setenv("JWT_SALT", "test-salt")
	os.Setenv("SinglePrice", "500")
	os.Setenv("PayExpireMinutes", "30")
	os.Exit(m.Run())
}

func newTestPaymentService(db *gorm.DB, expireMinutes int64) *PaymentService {
	return NewPaymentService(db, testutil.FakeGatewayInstance(), expireMinutes)
}

func mustCreatePay(t *testing.T, ps *PaymentService, db *gorm.DB, businessType, businessId, userId string, value int64) *models.Pay {
	t.Helper()
	pay, err := ps.CreatePay(context.Background(), db, businessType, businessId, userId, value)
	if err != nil {
		t.Fatalf("创建支付单失败: %s", err.Error())
	}
	return pay
}

func TestCreatePayIdempotentAndRebuild(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db, 30)
	ctx := context.Background()

	p1 := mustCreatePay(t, ps, db, "sign", "b1", "u1", 100)
	p2 := mustCreatePay(t, ps, db, "sign", "b1", "u1", 100)
	if p1.UUID != p2.UUID {
		t.Fatalf("幂等创建应返回同一支付单: %s vs %s", p1.UUID, p2.UUID)
	}
	if p1.ExpireTime-time.Now().Unix() != 30*60 {
		t.Errorf("过期时间应为 30 分钟后: %d", p1.ExpireTime)
	}

	// 过期后重建：旧单作废，生成新单
	if err := db.Model(&models.Pay{}).Where("uuid = ?", p1.UUID).Update("expire_time", time.Now().Unix()-1).Error; err != nil {
		t.Fatalf("修改过期时间失败: %s", err.Error())
	}
	p3 := mustCreatePay(t, ps, db, "sign", "b1", "u1", 100)
	if p3.UUID == p1.UUID {
		t.Fatal("过期后应新建支付单")
	}
	var oldPay models.Pay
	if err := db.First(&oldPay, "uuid = ?", p1.UUID).Error; err != nil {
		t.Fatalf("查询旧单失败: %s", err.Error())
	}
	if oldPay.Status != models.PayStatusVoid {
		t.Errorf("旧单应作废, got %d", oldPay.Status)
	}

	// 业务单已有已完成支付单时不允许再建
	if err := db.Model(&models.Pay{}).Where("uuid = ?", p3.UUID).Update("status", models.PayStatusPaid).Error; err != nil {
		t.Fatalf("更新状态失败: %s", err.Error())
	}
	if _, err := ps.CreatePay(ctx, db, "sign", "b1", "u1", 100); err == nil {
		t.Fatal("已存在完成支付单时应报错")
	}
}

func TestConfirmPaidAndBusinessLink(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db, 30)
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
	pay := mustCreatePay(t, ps, db, "sign", sign.UUID, "u1", 100)

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
	ps := newTestPaymentService(db, 30)
	ctx := context.Background()

	// 已过期
	expired := mustCreatePay(t, ps, db, "sign", "b2", "u1", 100)
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
	vo := mustCreatePay(t, ps, db, "sign", "b3", "u1", 100)
	if err := db.Model(&models.Pay{}).Where("uuid = ?", vo.UUID).Update("status", models.PayStatusVoid).Error; err != nil {
		t.Fatal(err)
	}
	if err := ps.ConfirmPaid(ctx, vo.UUID); err != ErrPayVoided {
		t.Errorf("作废单应报 ErrPayVoided, got %v", err)
	}

	// 已退款
	rf := mustCreatePay(t, ps, db, "sign", "b4", "u1", 100)
	if err := db.Model(&models.Pay{}).Where("uuid = ?", rf.UUID).Update("status", models.PayStatusRefunded).Error; err != nil {
		t.Fatal(err)
	}
	if err := ps.ConfirmPaid(ctx, rf.UUID); err != ErrPayRefunded {
		t.Errorf("退款单应报 ErrPayRefunded, got %v", err)
	}
}

func TestRefundRules(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db, 30)
	ctx := context.Background()

	pending := mustCreatePay(t, ps, db, "sign", "b9", "u1", 100)
	if err := ps.RefundPay(ctx, pending.UUID); err != ErrPayNotPaid {
		t.Errorf("待支付单退款应报 ErrPayNotPaid, got %v", err)
	}

	paid := mustCreatePay(t, ps, db, "sign", "b10", "u1", 100)
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

	vo := mustCreatePay(t, ps, db, "sign", "b11", "u1", 100)
	db.Model(&models.Pay{}).Where("uuid = ?", vo.UUID).Update("status", models.PayStatusVoid)
	if err := ps.RefundPay(ctx, vo.UUID); err != ErrPayVoided {
		t.Errorf("作废单退款应报 ErrPayVoided, got %v", err)
	}
}

func TestCancelMemberOrderVoidsPay(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db, 30)
	ms := NewMemberService(db, ps)
	ctx := context.Background()

	plan := models.MemberPlan{Name: "月卡", Type: "month", Value: 3000, Description: "一个月"}
	if err := gorm.G[models.MemberPlan](db).Create(ctx, &plan); err != nil {
		t.Fatal(err)
	}
	order := models.MemberOrders{PlanId: plan.UUID, UserId: "u1", Status: models.MemberOrderStatusPendingPay}
	if err := gorm.G[models.MemberOrders](db).Create(ctx, &order); err != nil {
		t.Fatal(err)
	}
	pay := mustCreatePay(t, ps, db, "member", order.UUID, "u1", 3000)

	n, err := ms.CancelMemberOrder(ctx, "u1", order.UUID)
	if err != nil || n == 0 {
		t.Fatalf("取消订单失败: %v, rows=%d", err, n)
	}
	var orderDb models.MemberOrders
	db.First(&orderDb, "uuid = ?", order.UUID)
	if orderDb.Status != models.MemberOrderStatusCanceled {
		t.Errorf("订单应取消, got %d", orderDb.Status)
	}
	var payDb models.Pay
	db.First(&payDb, "uuid = ?", pay.UUID)
	if payDb.Status != models.PayStatusVoid {
		t.Errorf("支付单应联动作废, got %d", payDb.Status)
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
	ps := newTestPaymentService(db, 30)
	ss := NewSignService(db, ps)
	ctx := context.Background()

	sign := models.Sign{UserId: "u1", Status: models.SignStatusActive, StartAt: time.Now().Unix() - 3600}
	if err := gorm.G[models.Sign](db).Create(ctx, &sign); err != nil {
		t.Fatal(err)
	}
	res, err := ss.FinishSignData(ctx, "u1", sign.UUID)
	if err != nil {
		t.Fatalf("结算失败: %s", err.Error())
	}
	if res.PayId == "" || res.PayUrl == "" || res.Value != 1000 {
		t.Errorf("结算结果异常: %+v", res)
	}
	var signDb models.Sign
	db.First(&signDb, "uuid = ?", sign.UUID)
	if signDb.Status != models.SignStatusPendingPay {
		t.Errorf("签到单应为待支付, got %d", signDb.Status)
	}
	var payDb models.Pay
	if err := db.Where("business_id = ?", sign.UUID).First(&payDb).Error; err != nil {
		t.Fatalf("支付单应存在: %s", err.Error())
	}

	// 已结算但未支付可幂等重试，返回同一支付单
	res2, err := ss.FinishSignData(ctx, "u1", sign.UUID)
	if err != nil {
		t.Fatalf("幂等重试结算失败: %s", err.Error())
	}
	if res2.PayId != res.PayId {
		t.Errorf("重复结算应返回同一支付单, got %s want %s", res2.PayId, res.PayId)
	}
}

func TestSignResettleAfterExpiry(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db, 30)
	ss := NewSignService(db, ps)
	ctx := context.Background()

	sign := models.Sign{UserId: "u1", Status: models.SignStatusActive, StartAt: time.Now().Unix() - 3600}
	if err := gorm.G[models.Sign](db).Create(ctx, &sign); err != nil {
		t.Fatal(err)
	}
	res1, err := ss.FinishSignData(ctx, "u1", sign.UUID)
	if err != nil {
		t.Fatalf("首次结算失败: %s", err.Error())
	}

	// 支付单过期后重新结算：旧单作废，生成新单
	if err := db.Model(&models.Pay{}).Where("uuid = ?", res1.PayId).Update("expire_time", time.Now().Unix()-1).Error; err != nil {
		t.Fatal(err)
	}
	res2, err := ss.FinishSignData(ctx, "u1", sign.UUID)
	if err != nil {
		t.Fatalf("过期后重新结算失败: %s", err.Error())
	}
	if res2.PayId == res1.PayId {
		t.Fatal("过期后应生成新支付单")
	}
	var oldPay models.Pay
	db.First(&oldPay, "uuid = ?", res1.PayId)
	if oldPay.Status != models.PayStatusVoid {
		t.Errorf("旧支付单应作废, got %d", oldPay.Status)
	}
	// 新支付单可正常确认
	if err := ps.ConfirmPaid(ctx, res2.PayId); err != nil {
		t.Fatalf("新支付单确认失败: %s", err.Error())
	}
}

func TestMemberFinishFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db, 30)
	ms := NewMemberService(db, ps)
	ctx := context.Background()

	plan := models.MemberPlan{Name: "周卡", Type: "week", Value: 2000, Description: "一周"}
	if err := gorm.G[models.MemberPlan](db).Create(ctx, &plan); err != nil {
		t.Fatal(err)
	}
	order := models.MemberOrders{PlanId: plan.UUID, UserId: "u1", Status: models.MemberOrderStatusCreated}
	if err := gorm.G[models.MemberOrders](db).Create(ctx, &order); err != nil {
		t.Fatal(err)
	}
	res, err := ms.FinishMemberOrder(ctx, "u1", order.UUID)
	if err != nil {
		t.Fatalf("会员结算失败: %s", err.Error())
	}
	if res.PayId == "" || res.Value != 2000 {
		t.Errorf("结算结果异常: %+v", res)
	}
	var orderDb models.MemberOrders
	db.First(&orderDb, "uuid = ?", order.UUID)
	if orderDb.Status != models.MemberOrderStatusPendingPay {
		t.Errorf("订单应为待支付, got %d", orderDb.Status)
	}
	// 确认支付后订单联动为已支付
	if err := ps.ConfirmPaid(ctx, res.PayId); err != nil {
		t.Fatalf("确认支付失败: %s", err.Error())
	}
	db.First(&orderDb, "uuid = ?", order.UUID)
	if orderDb.Status != models.MemberOrderStatusPaid {
		t.Errorf("订单应联动为已支付, got %d", orderDb.Status)
	}
}

func TestMemberResettleAfterExpiry(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db, 30)
	ms := NewMemberService(db, ps)
	ctx := context.Background()

	plan := models.MemberPlan{Name: "周卡", Type: "week", Value: 2000, Description: "一周"}
	if err := gorm.G[models.MemberPlan](db).Create(ctx, &plan); err != nil {
		t.Fatal(err)
	}
	order := models.MemberOrders{PlanId: plan.UUID, UserId: "u1", Status: models.MemberOrderStatusCreated}
	if err := gorm.G[models.MemberOrders](db).Create(ctx, &order); err != nil {
		t.Fatal(err)
	}
	res1, err := ms.FinishMemberOrder(ctx, "u1", order.UUID)
	if err != nil {
		t.Fatalf("首次结算失败: %s", err.Error())
	}
	if err := db.Model(&models.Pay{}).Where("uuid = ?", res1.PayId).Update("expire_time", time.Now().Unix()-1).Error; err != nil {
		t.Fatal(err)
	}
	res2, err := ms.FinishMemberOrder(ctx, "u1", order.UUID)
	if err != nil {
		t.Fatalf("过期后重新结算失败: %s", err.Error())
	}
	if res2.PayId == res1.PayId {
		t.Fatal("过期后应生成新支付单")
	}
	var oldPay models.Pay
	db.First(&oldPay, "uuid = ?", res1.PayId)
	if oldPay.Status != models.PayStatusVoid {
		t.Errorf("旧支付单应作废, got %d", oldPay.Status)
	}
}

func TestExpireSweeper(t *testing.T) {
	db := testutil.SetupTestDB(t)
	ps := newTestPaymentService(db, 30)
	ctx := context.Background()

	pay := mustCreatePay(t, ps, db, "sign", "b12", "u1", 100)
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

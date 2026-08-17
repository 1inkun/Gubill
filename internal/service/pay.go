package service

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/payment"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

var (
	ErrPayNoExist       = models.NewBusinessError(400, "支付单不存在")
	ErrPayAlreadyPaid   = models.NewBusinessError(400, "支付单已完成支付")
	ErrPayExpired       = models.NewBusinessError(400, "支付单已过期")
	ErrPayVoided        = models.NewBusinessError(400, "支付单已作废")
	ErrPayRefunded      = models.NewBusinessError(400, "支付单已退款")
	ErrPayNotPaid       = models.NewBusinessError(400, "支付单未支付，无法退款")
	ErrPayOwnerMismatch = models.NewBusinessError(400, "无权查看该支付单")
)

// PaymentService 负责支付单的确认、退款与查询。
// gateway 为真实支付渠道预留位（TODO(支付接入)），未接入时为 nil。
type PaymentService struct {
	db      *gorm.DB
	gateway payment.Gateway
}

// NewPaymentService 创建支付服务。
func NewPaymentService(db *gorm.DB, gateway payment.Gateway) *PaymentService {
	return &PaymentService{db: db, gateway: gateway}
}

// ConfirmPaid 将待支付单流转为已支付，并联动更新对应业务单状态。
// 所需数据（金额、业务类型、业务单号）全部从 pays 读取。
func (s *PaymentService) ConfirmPaid(ctx context.Context, payId string) error {
	pay, err := s.loadPay(ctx, payId)
	if err != nil {
		return err
	}
	switch pay.Status {
	case models.PayStatusPaid:
		return ErrPayAlreadyPaid
	case models.PayStatusRefunded:
		return ErrPayRefunded
	case models.PayStatusVoid:
		return ErrPayVoided
	case models.PayStatusPending:
		if pay.ExpireTime <= time.Now().Unix() {
			if _, err := gorm.G[models.Pay](s.db).Where("uuid = ?", payId).Update(ctx, "status", models.PayStatusVoid); err != nil {
				return models.NewDatabaseErr(err)
			}
			return ErrPayExpired
		}
	default:
		return ErrPayVoided
	}

	now := time.Now().Unix()
	txnId := pay.TransactionId
	if txnId == "" {
		if s.gateway != nil {
			txnId = s.gateway.NewTransactionId()
		} else {
			// 网关未配置时（线下/现金收款），使用本地手工流水号；
			// 接入真实渠道后应优先使用回调携带的渠道交易号。
			txnId = "manual-" + pay.UUID
		}
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 双检：避免并发重复确认
		check, e := gorm.G[models.Pay](tx).Where("uuid = ?", payId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(check) == 0 {
			return ErrPayNoExist
		}
		if check[0].Status != models.PayStatusPending {
			switch check[0].Status {
			case models.PayStatusVoid:
				return ErrPayVoided
			case models.PayStatusRefunded:
				return ErrPayRefunded
			default:
				return ErrPayAlreadyPaid
			}
		}
		if e := tx.Model(&models.Pay{}).Where("uuid = ?", payId).Updates(map[string]any{
			"status":         models.PayStatusPaid,
			"pay_at":         now,
			"transaction_id": txnId,
		}).Error; e != nil {
			return models.NewDatabaseErr(e)
		}
		// 联动业务单
		switch pay.BusinessType {
		case models.PayBusinessSign:
			if _, e := gorm.G[models.Sign](tx).Where("uuid = ? AND status = ?", pay.BusinessId, models.SignStatusPendingPay).Update(ctx, "status", models.SignStatusFinished); e != nil {
				return models.NewDatabaseErr(e)
			}
		case models.PayBusinessMember:
			if _, e := gorm.G[models.MemberOrders](tx).Where("uuid = ? AND status = ?", pay.BusinessId, models.MemberOrderStatusPendingPay).Update(ctx, "status", models.MemberOrderStatusPaid); e != nil {
				return models.NewDatabaseErr(e)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// RefundPay 全额退款：仅已支付单可退，退款记录在支付单上，不回滚业务单状态。
func (s *PaymentService) RefundPay(ctx context.Context, payId string) error {
	pay, err := s.loadPay(ctx, payId)
	if err != nil {
		return err
	}
	if pay.Status != models.PayStatusPaid {
		switch pay.Status {
		case models.PayStatusPending:
			return ErrPayNotPaid
		case models.PayStatusRefunded:
			return ErrPayRefunded
		case models.PayStatusVoid:
			return ErrPayVoided
		default:
			return ErrPayNotPaid
		}
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		check, e := gorm.G[models.Pay](tx).Where("uuid = ?", payId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(check) == 0 {
			return ErrPayNoExist
		}
		if check[0].Status != models.PayStatusPaid {
			return ErrPayNotPaid
		}
		// 网关未配置时跳过渠道退款调用（线下退款场景）；
		// 接入真实渠道后此处会调用渠道退款接口。
		if s.gateway != nil {
			if e := s.gateway.Refund(ctx, &payment.RefundRequest{
				PayId:  payId,
				Value:  pay.Value,
				Refund: pay.Value,
			}); e != nil {
				return models.NewInternalError(500, "网关退款失败", e)
			}
		}
		if e := tx.Model(&models.Pay{}).Where("uuid = ?", payId).Updates(map[string]any{
			"status":       models.PayStatusRefunded,
			"refund_value": pay.Value,
			"refund_at":    time.Now().Unix(),
		}).Error; e != nil {
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// ListUserPays 查询用户自己的支付记录（分页）。
func (s *PaymentService) ListUserPays(ctx context.Context, userId string, page, pageSize int) ([]models.PayRes, error) {
	var datas []models.PayRes
	tx := s.db.Model(&models.Pay{}).Where("user_id = ?", userId).
		Select("uuid, user_id, business_type, business_id, channel, value, status, transaction_id, expire_time, pay_at, refund_value, refund_at, created_at").
		Scopes(utils.Paginate(page, pageSize)).Find(&datas)
	if tx.Error != nil {
		return nil, models.NewDatabaseErr(tx.Error)
	}
	if len(datas) == 0 {
		return nil, nil
	}
	return datas, nil
}

// GetUserPay 查询用户自己的单笔支付单。
func (s *PaymentService) GetUserPay(ctx context.Context, userId, payId string) (models.PayRes, error) {
	pay, err := s.loadPay(ctx, payId)
	if err != nil {
		return models.PayRes{}, err
	}
	if pay.UserId != userId {
		return models.PayRes{}, ErrPayOwnerMismatch
	}
	return pay.ToPayRes(), nil
}

// ListAllPays 管理端查询全部支付记录（分页）。
func (s *PaymentService) ListAllPays(ctx context.Context, page, pageSize int) ([]models.PayRes, error) {
	var datas []models.PayRes
	tx := s.db.Model(&models.Pay{}).
		Select("uuid, user_id, business_type, business_id, channel, value, status, transaction_id, expire_time, pay_at, refund_value, refund_at, created_at").
		Scopes(utils.Paginate(page, pageSize)).Find(&datas)
	if tx.Error != nil {
		return nil, models.NewDatabaseErr(tx.Error)
	}
	if len(datas) == 0 {
		return nil, nil
	}
	return datas, nil
}

// GetPay 管理端查询单笔支付单。
func (s *PaymentService) GetPay(ctx context.Context, payId string) (models.PayRes, error) {
	pay, err := s.loadPay(ctx, payId)
	if err != nil {
		return models.PayRes{}, err
	}
	return pay.ToPayRes(), nil
}

// GetPayByBusinessId 按业务单号查询支付单（管理端联动使用）。
func (s *PaymentService) GetPayByBusinessId(ctx context.Context, businessId string) (models.PayRes, error) {
	datas, err := gorm.G[models.Pay](s.db).Where("business_id = ?", businessId).Order("created_at DESC").Limit(1).Find(ctx)
	if err != nil {
		return models.PayRes{}, models.NewDatabaseErr(err)
	}
	if len(datas) == 0 {
		return models.PayRes{}, ErrPayNoExist
	}
	return datas[0].ToPayRes(), nil
}

// TodayStats 返回今日已支付笔数与金额合计。
func (s *PaymentService) TodayStats(ctx context.Context) (count int64, sum int64, err error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	var row struct {
		Cnt int64
		Sum int64
	}
	tx := s.db.Model(&models.Pay{}).
		Select("COUNT(*) AS cnt, COALESCE(SUM(value), 0) AS sum").
		Where("status = ? AND pay_at >= ?", models.PayStatusPaid, start).
		Scan(&row)
	if tx.Error != nil {
		return 0, 0, models.NewDatabaseErr(tx.Error)
	}
	return row.Cnt, row.Sum, nil
}

// PendingCount 返回待支付单数量。
func (s *PaymentService) PendingCount(ctx context.Context) (int64, error) {
	var count int64
	tx := s.db.Model(&models.Pay{}).Where("status = ?", models.PayStatusPending).Count(&count)
	if tx.Error != nil {
		return 0, models.NewDatabaseErr(tx.Error)
	}
	return count, nil
}

// StartExpireSweeper 启动后台定时任务：定期作废过期未支付单。
func (s *PaymentService) StartExpireSweeper(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n, err := s.expirePending(ctx)
			if err != nil {
				log.Printf("过期支付单清理失败:%s", err.Error())
				continue
			}
			if n > 0 {
				log.Printf("已作废 %d 笔过期支付单", n)
			}
		}
	}
}

func (s *PaymentService) expirePending(ctx context.Context) (int64, error) {
	n, err := gorm.G[models.Pay](s.db).
		Where("status = ? AND expire_time < ?", models.PayStatusPending, time.Now().Unix()).
		Update(ctx, "status", models.PayStatusVoid)
	if err != nil {
		return 0, models.NewDatabaseErr(err)
	}
	return int64(n), nil
}

func (s *PaymentService) loadPay(ctx context.Context, payId string) (models.Pay, error) {
	datas, err := gorm.G[models.Pay](s.db).Where("uuid = ?", payId).Limit(1).Find(ctx)
	if err != nil {
		return models.Pay{}, models.NewDatabaseErr(err)
	}
	if len(datas) == 0 {
		return models.Pay{}, ErrPayNoExist
	}
	return datas[0], nil
}

// SinglePrice 读取签到单价（单位：分），默认 500。
func SinglePrice() int64 {
	v := os.Getenv("SinglePrice")
	if v == "" {
		return 500
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n <= 0 {
		return 500
	}
	return n
}

// PayExpireSeconds 读取支付单有效期（秒），默认 30 分钟。
func PayExpireSeconds() int64 {
	v := os.Getenv("PayExpireMinutes")
	if v == "" {
		return 30 * 60
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n <= 0 {
		return 30 * 60
	}
	return n * 60
}

package repository

import (
	"context"
	"log"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

type SignQuery struct {
	db *gorm.DB
}

func NewSignQuery(db *gorm.DB) *SignQuery {
	return &SignQuery{db: db}
}

var (
	ErrExistSignData   = models.NewBusinessError(400, "存在未结算的签到订单")
	ErrUserNoExist     = models.NewBusinessError(400, "用户不存在")
	ErrSignDataNoExist = models.NewBusinessError(400, "签到订单不存在")
	ErrWrongUser       = models.NewBusinessError(400, "结算的用户有误")
	ErrAlreadyFinished = models.NewBusinessError(400, "订单已经完成")
	ErrExistPayData    = models.NewBusinessError(400, "存在尚未支付的订单")
)

// 面向用户的接口

func (q *SignQuery) GenerateNewSignData(ctx context.Context, userId string) (string, error) {
	var newData = models.Sign{
		UserId:  userId,
		Status:  0,
		StartAt: time.Now().Unix(),
	}
	// 在事务中处理
	err := q.db.Transaction(func(tx *gorm.DB) error {
		// 检查是否存在未结算的订单
		checkSign, e := gorm.G[models.Sign](tx).Select("uuid").Where("user_id = ? AND status != 1", userId).Find(ctx)
		if e != nil {
			log.Println(e.Error())
			return models.NewDatabaseErr(e)
		}
		if len(checkSign) > 0 {
			return ErrExistSignData
		}
		checkPay, e := gorm.G[models.Pay](tx).Select("uuid").Where("user_id = ? AND status != 1", userId).Find(ctx)
		if e != nil {
			log.Println(e.Error())
			return models.NewDatabaseErr(e)
		}
		if len(checkPay) > 0 {
			return ErrExistPayData
		}
		// 创建新订单
		e = gorm.G[models.Sign](tx).Create(ctx, &newData)
		if e != nil {
			log.Println(e.Error())
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return newData.UUID, nil
}

func (q *SignQuery) GetUserSignData(ctx context.Context, userId string, status int64, page int, pageSize int) ([]models.SignRes, error) {
	var results []models.SignRes
	tx := q.db.Model(models.Sign{}).WithContext(ctx).Where("user_id = ? AND status = ?", userId, status).Scopes(utils.Paginate(page, pageSize)).Find(&results, &models.SignRes{})
	if tx.Error != nil {
		log.Println(tx.Error.Error())
		return nil, models.NewDatabaseErr(tx.Error)
	}
	return results, nil
}

func (q *SignQuery) GetSignDataById(ctx context.Context, signId string) (models.SignRes, error) {
	results, err := gorm.G[models.Sign](q.db).Where("uuid = ?", signId).Find(ctx)
	if err != nil {
		log.Println(err.Error())
		return models.SignRes{}, models.NewDatabaseErr(err)
	}
	if len(results) == 0 {
		return models.SignRes{}, ErrSignDataNoExist
	}
	// 处理data
	data := results[0]
	result := models.SignRes{
		UUID:    data.UUID,
		UserId:  data.UserId,
		StartAt: data.StartAt,
		EndAt:   data.EndAt,
		Status:  data.Status,
		Value:   data.Value,
	}
	return result, nil
}

func (s *SignQuery) FinishSignData(ctx context.Context, userId string, signId string) (string, error) {
	var newPay = models.Pay{
		BusinessType: "sign",
		Status:       0,
		UserId:       userId,
		ExpireTime:   time.Now().Add(30 * time.Minute).Unix(),
	}
	now := time.Now().Unix()
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var signData models.Sign
		signData, e := gorm.G[models.Sign](tx).Where("uuid = ? AND status = 0", signId).Limit(1).First(ctx)
		if e != nil {
			log.Println(e.Error())
			if e == gorm.ErrRecordNotFound {
				return ErrSignDataNoExist
			}
			return models.NewDatabaseErr(e)
		}
		if signData.UserId != userId {
			return ErrWrongUser
		}
		duration := now - signData.StartAt
		// 500为单价，单位为分，暂时先用来占位
		totalPrice := utils.CalculateTotalPrice(duration, 500)
		// 构造支付订单数据,构造完成后插入数据
		newPay.BusinessId = signData.UUID
		newPay.Value = int64(totalPrice)
		// 插入数据
		e = gorm.G[models.Pay](tx).Create(ctx, &newPay)
		if e != nil {
			log.Println(e.Error())
			return models.NewDatabaseErr(e)
		}
		// 更新订单表中的数据
		var updateSign models.Sign
		// 更新为2表示待支付
		updateSign.Status = 2
		updateSign.EndAt = now
		updateSign.Value = int64(totalPrice)
		_, e = gorm.G[models.Sign](tx).Where("uuid = ?", signId).Updates(ctx, updateSign)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return newPay.UUID, nil
}

// 面向管理的接口
func (s *SignQuery) GetAllSignData(ctx context.Context, page int, pageSize int) ([]models.SignRes, error) {
	var results []models.SignRes
	tx := s.db.Model(&models.Sign{}).WithContext(ctx).Scopes(utils.Paginate(page, pageSize)).Find(&results, &models.SignRes{})
	if tx.Error != nil {
		log.Println(tx.Error.Error())
		return nil, models.NewDatabaseErr(tx.Error)
	}
	return results, nil
}

func (s *SignQuery) UpdateSignData(ctx context.Context, signId string, status int64, value int64) (int, error) {
	var rows int
	// 构造更新内容
	var newData models.Sign
	newData.Status = status
	if newData.Status != 0 {
		newData.EndAt = time.Now().Unix()
		newData.Value = value
	} else {
		newData.Value = value
	}
	// 在事务中处理
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 检查当前签到数据是否存在
		_, e := gorm.G[models.Sign](tx).Select("uuid").Where("uuid = ?", signId).First(ctx)
		if e != nil {
			log.Println(e)
			if e == gorm.ErrRecordNotFound {
				return ErrSignDataNoExist
			}
			return models.NewDatabaseErr(e)
		}
		// 写入修改
		rows, e = gorm.G[models.Sign](tx).Where("uuid = ?", signId).Updates(ctx, newData)
		if e != nil {
			log.Println(e)
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return rows, nil
}

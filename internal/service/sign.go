package service

import (
	"context"
	"log"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

type SignService struct {
	db *gorm.DB
}

func NewSignService(db *gorm.DB) *SignService {
	return &SignService{
		db: db,
	}
}

var (
	ErrExistSignData   = models.NewBusinessError(400, "存在未结算的签到订单")
	ErrUserNoExist     = models.NewBusinessError(400, "用户不存在")
	ErrSignDataNoExist = models.NewBusinessError(400, "签到订单不存在")
	ErrWrongUser       = models.NewBusinessError(400, "结算的用户有误")
	ErrAlreadyFinished = models.NewBusinessError(400, "订单已经完成")
)

func (s *SignService) GenerateNewSignData(ctx context.Context, userId string) (string, error) {
	var newData models.Sign
	newData.UserId = userId
	newData.Status = 0
	newData.StartAt = time.Now().Unix()
	// 在事务中处理
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 查表确定用户真实存在
		userData, e := gorm.G[models.User](tx).Where("uuid = ?", userId).Find(ctx)
		if e != nil {
			log.Println(e.Error())
			return ErrDatabaseErr
		}
		if len(userData) == 0 {
			return ErrUserNoExist
		}
		// 确定用户存在后继续查表,确保没有正在进行中的签到
		check, e := gorm.G[models.Sign](tx).Where("user_id = ? AND status = 0", userId).Find(ctx)
		if e != nil {
			return ErrDatabaseErr
		}
		if len(check) != 0 {
			return ErrExistSignData
		}
		// 确认无误后添加新的签到数据
		e = gorm.G[models.Sign](tx).Create(ctx, &newData)
		if e != nil {
			return ErrDatabaseErr
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return newData.UUID, nil
}

func (s *SignService) GetUserSignData(ctx context.Context, userId string, status int64) ([]models.SignRes, error) {
	var signDatas []models.Sign
	signDatas, err := gorm.G[models.Sign](s.db).Where("user_id = ? AND status = ?", userId, status).Find(ctx)
	if err != nil {
		return nil, ErrDatabaseErr
	}
	if len(signDatas) == 0 {
		return nil, nil
	}
	var resData []models.SignRes
	for _, v := range signDatas {
		var temp = models.SignRes{
			UUID:    v.UUID,
			UserId:  userId,
			StartAt: v.StartAt,
			EndAt:   v.EndAt,
			Status:  v.Status,
			Value:   v.Value,
		}
		resData = append(resData, temp)
	}

	return resData, nil
}

func (s *SignService) GetSignData(ctx context.Context, userId string, signId string) (models.SignRes, error) {
	var signData []models.Sign
	signData, err := gorm.G[models.Sign](s.db).Where("uuid = ?", signId).Find(ctx)
	if err != nil {
		return models.SignRes{}, ErrDatabaseErr
	}
	if len(signData) == 0 {
		return models.SignRes{}, nil
	}
	data := signData[0]
	if data.UserId != userId {
		return models.SignRes{}, ErrWrongUser
	}
	signRes := models.SignRes{
		UUID:    data.UUID,
		UserId:  data.UserId,
		StartAt: data.StartAt,
		EndAt:   data.EndAt,
		Status:  data.Status,
		Value:   data.Value,
	}
	return signRes, nil
}

func (s *SignService) FinishSignData(ctx context.Context, userId string, signId string) (string, error) {
	var newPay models.Pay
	newPay.BusinessType = "sign"
	newPay.Status = 0
	now := time.Now().Unix()
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var signData []models.Sign
		signData, e := gorm.G[models.Sign](tx).Where("uuid = ?", signId).Limit(1).Find(ctx)
		if e != nil {
			log.Println(e.Error())
			return ErrDatabaseErr
		}
		if len(signData) == 0 {
			return ErrSignDataNoExist
		}
		data := signData[0]
		if data.Status == 1 {
			return ErrAlreadyFinished
		}
		if data.UserId != userId {
			return ErrWrongUser
		}
		duration := now - data.StartAt
		totalPrice := utils.CalculateTotalPrice(duration, 500)
		// 构造支付订单数据,构造完成后插入数据
		newPay.BusinessId = data.UUID
		newPay.Value = int64(totalPrice)
		newPay.ExpireTime = time.Now().Add(30 * time.Minute).Unix()
		e = gorm.G[models.Pay](tx).Create(ctx, &newPay)
		if e != nil {
			log.Println(e.Error())
			return ErrDatabaseErr
		}
		// 更新订单表中的数据
		var updateSign models.Sign
		updateSign.Status = 2
		updateSign.EndAt = now
		updateSign.Value = int64(totalPrice)
		_, e = gorm.G[models.Sign](tx).Where("uuid = ?", signId).Updates(ctx, updateSign)
		if e != nil {
			log.Println(e)
			return ErrDatabaseErr
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return newPay.UUID, nil
}

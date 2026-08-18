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
	ErrExistPayData    = models.NewBusinessError(400, "存在尚未支付的订单")
)

func (s *SignService) GenerateNewSignData(ctx context.Context, userId string) (string, error) {
	var newData models.Sign
	newData.UserId = userId
	newData.Status = 0
	newData.StartAt = time.Now().Unix()
	// 在事务中处理
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 查表确定用户真实存在
		userData, e := gorm.G[models.User](tx).Select("uuid").Where("uuid = ?", userId).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(userData) == 0 {
			return ErrUserNoExist
		}
		// 确定用户存在后继续查表,确保该用户没有正在进行中的签到
		check, e := gorm.G[models.Sign](tx).Select("uuid").Where("user_id = ? AND status = 0", userId).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(check) != 0 {
			return ErrExistSignData
		}
		// 确定用户存在后继续查表,确保用户没有暂未支付的订单
		payCheck, e := gorm.G[models.Pay](tx).Select("uuid").Where("user_id = ? AND status = 0", userId).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(payCheck) != 0 {
			return ErrExistPayData
		}
		// 确认无误后添加新的签到数据
		e = gorm.G[models.Sign](tx).Create(ctx, &newData)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return newData.UUID, nil
}

func (s *SignService) GetAllSignData(ctx context.Context, page int, pageSize int) ([]models.SignRes, error) {
	var datas []models.SignRes
	tx := s.db.Model(&models.Sign{}).Select(&models.SignRes{}).WithContext(ctx).Scopes(utils.Paginate(page, pageSize)).Find(&datas)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
		return nil, models.NewDatabaseErr(tx.Error)
	}
	if len(datas) == 0 {
		return nil, nil
	}
	return datas, nil
}

func (s *SignService) GetSignDataBySignId(ctx context.Context, signId string) (models.SignRes, error) {
	var rawData []models.Sign
	rawData, err := gorm.G[models.Sign](s.db).Where("uuid = ?", signId).Limit(1).Find(ctx)
	if err != nil {
		log.Println(err)
		return models.SignRes{}, models.NewDatabaseErr(err)
	}
	if len(rawData) == 0 {
		return models.SignRes{}, ErrNoSuchData
	}
	data := rawData[0]
	var resData = models.SignRes{
		UUID:    data.UUID,
		UserId:  data.UserId,
		StartAt: data.StartAt,
		EndAt:   data.EndAt,
		Status:  data.Status,
		Value:   data.Value,
	}
	return resData, nil
}

func (s *SignService) GetUserSignData(ctx context.Context, userId string, status int64) ([]models.SignRes, error) {
	var signDatas []models.Sign
	signDatas, err := gorm.G[models.Sign](s.db).Where("user_id = ? AND status = ?", userId, status).Find(ctx)
	if err != nil {
		return nil, models.NewDatabaseErr(err)
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
		return models.SignRes{}, models.NewDatabaseErr(err)
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

func (s *SignService) UpdateSignData(ctx context.Context, signId string, status int64, value int64) (int, error) {
	var res int
	var newData models.Sign
	newData.Status = status
	if newData.Status != 0 {
		newData.EndAt = time.Now().Unix()
		newData.Value = value
	} else {
		newData.Value = value
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 检查是否存在对应的数据
		var check []models.Sign
		check, e := gorm.G[models.Sign](tx).Select("uuid").Where("uuid = ?", signId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(check) == 0 {
			return ErrNoSuchData
		}
		// 实际更新
		res, e = gorm.G[models.Sign](tx).Where("uuid = ?", signId).Updates(ctx, newData)
		return nil
	})
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *SignService) FinishSignData(ctx context.Context, userId string, signId string) (string, error) {
	var newPay models.Pay
	newPay.BusinessType = "sign"
	newPay.Status = 0
	now := time.Now().Unix()
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var signData []models.Sign
		signData, e := gorm.G[models.Sign](tx).Where("uuid = ? AND status = 0", signId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(signData) == 0 {
			return ErrSignDataNoExist
		}
		data := signData[0]
		if data.UserId != userId {
			return ErrWrongUser
		}
		duration := now - data.StartAt
		totalPrice := utils.CalculateTotalPrice(duration, 500)
		// 构造支付订单数据,构造完成后插入数据
		newPay.UserId = userId
		newPay.BusinessId = data.UUID
		newPay.Value = int64(totalPrice)
		newPay.ExpireTime = time.Now().Add(30 * time.Minute).Unix()
		// 插入数据
		e = gorm.G[models.Pay](tx).Create(ctx, &newPay)
		if e != nil {
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

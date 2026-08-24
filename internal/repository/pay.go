package repository

import (
	"context"
	"log"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

type PayQuery struct {
	db *gorm.DB
}

func NewPayQuery(db *gorm.DB) *PayQuery {
	return &PayQuery{db: db}
}

func (q *PayQuery) GetUserPayOrders(ctx context.Context, userId string, page int, pageSize int) ([]models.PayRes, error) {
	var results []models.PayRes
	tx := q.db.Model(&models.Pay{}).WithContext(ctx).Scopes(utils.Paginate(page, pageSize)).Where("user_id = ?", userId).Find(&results, &models.PayRes{})
	if tx.Error != nil {
		return nil, models.NewDatabaseErr(tx.Error)
	}
	return results, nil
}

func (q *PayQuery) GetUserPayOrderDataNums(ctx context.Context, userId string) (int64, error) {
	var rowDatas []models.Pay
	rowDatas, err := gorm.G[models.Pay](q.db).Select("uuid").Where("user_id", userId).Find(ctx)
	if err != nil {
		return 0, models.NewDatabaseErr(err)
	}
	nums := len(rowDatas)
	return int64(nums), err
}

func (q *PayQuery) GetPayOrdersById(ctx context.Context, payId string) (models.PayRes, error) {
	data, err := gorm.G[models.Pay](q.db).Where("uuid = ?", payId).First(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.PayRes{}, ErrExistPayData
		}
		log.Println(err)
		return models.PayRes{}, models.NewDatabaseErr(err)
	}
	var res = models.PayRes{
		UUID:         data.UUID,
		BusinessType: data.BusinessType,
		BusinessId:   data.BusinessId,
		UserId:       data.UserId,
		Value:        data.Value,
		Status:       data.Status,
		ExpireTime:   data.ExpireTime,
		PayAt:        data.PayAt,
	}
	return res, nil
}

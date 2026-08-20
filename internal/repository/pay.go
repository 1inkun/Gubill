package repository

import (
	"context"

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

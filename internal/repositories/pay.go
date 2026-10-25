package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/1inkun/billApp/internal/models"
	"github.com/redis/go-redis/v9"
)

type PayQuery struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewPayQuery(db *sql.DB, rdb *redis.Client) *PayQuery {
	return &PayQuery{
		db:  db,
		rdb: rdb,
	}
}

func (q *PayQuery) GetPayOrderDataById(ctx context.Context, orderId int64) (models.PayOrderData, error) {
	var data models.PayOrderData
	query := "SELECT `id`, `user_id`, `type`, `order_no`, `status`, `total_fee`, `description` FROM `pay_orders` WHERE `id` = ?;"
	row := q.db.QueryRowContext(ctx, query, orderId)
	if row.Err() != nil {
		log.Printf("在GetPayOrderDataById中查询数据库时出错:%s", row.Err())
		return models.PayOrderData{}, row.Err()
	}
	err := row.Scan(&data.Id, &data.UserId, &data.Type, &data.OrderNo, &data.Status, &data.TotalFee, &data.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.PayOrderData{}, models.ErrNoSuchPayOrder
		}
	}
	return data, nil
}
func (q *PayQuery) GetPayOrderDataByIdFromCache(ctx context.Context, key string) (models.PayOrderData, error) {
	var data models.PayOrderData
	val, err := q.rdb.Get(ctx, key).Result()
	if err != nil {
		return models.PayOrderData{}, err
	}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		log.Printf("在GetPayOrderDataByIdFromCache中反序列化数据出错:%s", err)
		return models.PayOrderData{}, err
	}
	return data, nil
}
func (q *PayQuery) GetPayOrderDataByUserId(ctx context.Context, userId int64, status, offset int64, limit int64) ([]models.PayOrderData, error) {
	var data []models.PayOrderData
	var err error
	var rows *sql.Rows
	if status < -1 {
		query := "SELECT `id`, `user_id`, `type`, `order_no`, `status`, `total_fee`, `description` FROM `pay_orders` WHERE `user_id` = ? LIMIT ? OFFSET ?;"
		rows, err = q.db.QueryContext(ctx, query, userId, limit, offset)
	} else {
		query := "SELECT `id`, `user_id`, `type`, `order_no`, `status`, `total_fee`, `description` FROM `pay_orders` WHERE `user_id` = ? AND `status` = ? LIMIT ? OFFSET ?;"
		rows, err = q.db.QueryContext(ctx, query, userId, status, limit, offset)
	}
	if err != nil {
		log.Printf("在GetPayOrderDataByUserId中查询数据时出错:%s", err)
		return nil, err
	}
	if rows.Err() != nil {
		log.Printf("在GetPayOrderDataByUserId中查询数据时出错:%s", rows.Err())
		return nil, rows.Err()
	}
	defer rows.Close()
	for rows.Next() {
		var i models.PayOrderData
		err := rows.Scan(&i.Id, &i.UserId, &i.Type, &i.OrderNo, &i.Status, &i.TotalFee, &i.Description)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			log.Printf("在GetPayOrderDataByUserId中查询数据时出错:%s", err)
			return nil, err
		}
		data = append(data, i)
	}
	return data, nil
}
func (q *PayQuery) SetCache(ctx context.Context, key string, value string) error {
	_, err := q.rdb.Set(ctx, key, value, 5*time.Minute).Result()
	if err != nil {
		return err
	}
	return nil
}

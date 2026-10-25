package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	"github.com/1inkun/billApp/internal/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type GoodQuery struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewGoodQuery(db *sql.DB, rdb *redis.Client) *GoodQuery {
	return &GoodQuery{
		db:  db,
		rdb: rdb,
	}
}

func (q *GoodQuery) GetGoods(ctx context.Context) ([]models.Good, error) {
	var data []models.Good
	query := "SELECT `id`, `name`, `type`, `description`, `status`, `price` FROM `goods`;"
	rows, _ := q.db.QueryContext(ctx, query)
	if rows.Err() != nil {
		log.Printf("在GetGoods中查询数据库时出现错误:%s", rows.Err())
		return nil, rows.Err()
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Good
		if err := rows.Scan(&i.Id, &i.Name, &i.Type, &i.Description, &i.Status, &i.Price); err != nil {
			if err == sql.ErrNoRows {
				// 没找到数据
				return nil, nil
			}
			log.Printf("在GetGoods中查询数据库时出现错误:%s", err)
			return nil, err
		}
		data = append(data, i)
	}
	return data, nil
}

func (q *GoodQuery) GetGoodsById(ctx context.Context, goodId int64) (models.Good, error) {
	var data models.Good
	query := "SELECT `id`, `name`, `type`, `description`, `status`, `price` FROM `goods` WHERE `id` = ?;"
	row := q.db.QueryRowContext(ctx, query, goodId)
	if row.Err() != nil {
		log.Printf("在GetGoodsById中查找数据库时出错:%s", row.Err())
		return models.Good{}, row.Err()
	}
	err := row.Scan(&data.Id, &data.Name, &data.Type, &data.Description, &data.Status, &data.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Good{}, nil
		}
		log.Printf("在GetGoodsById中查找数据库时出错:%s", err)
		return models.Good{}, err
	}
	return data, nil
}

func (q *GoodQuery) GetGoodsFromCache(ctx context.Context) ([]models.Good, error) {
	key := "goodsList"
	check, err := q.rdb.SCard(ctx, key).Result()
	if check == 0 || err != nil {
		// 没找到数据
		return nil, nil
	}
	val, err := q.rdb.SMembers(ctx, key).Result()
	if err != nil {
		log.Printf("在GetGoods中查询缓存时出现错误:%s", err)
		return nil, err
	}
	// 将查询到的字符串反序列化
	var data []models.Good
	var eg errgroup.Group
	var mu sync.Mutex
	eg.Go(func() error {
		for _, v := range val {
			mu.Lock()
			var i models.Good
			err := json.Unmarshal([]byte(v), &i)
			if err != nil {
				log.Printf("在GetGoods中反序列化缓存中的数据时出错:%s", err)
				return err
			}
			data = append(data, i)
			mu.Unlock()
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return data, nil
}

func (q *GoodQuery) SetCache(ctx context.Context, key string, value string) error {
	_, err := q.rdb.SAdd(ctx, key, value).Result()
	if err != nil {
		return err
	}
	return nil
}

func (q *GoodQuery) BuyGood(ctx context.Context, goodData models.Good, userId int64) (models.BuyGoodResponse, error) {
	query := "INSERT INTO pay_orders (`user_id`, `type`, `order_no`, `status`, `description`, `total_fee`) VALUES (?, ?, ?, ?, ?, ?);"
	res, err := q.db.ExecContext(ctx, query, userId, "goods", -1, 0, "会员商品购买", goodData.Price)
	if err != nil {
		log.Printf("在BuyGood中创建新订单失败:%s", err)
		return models.BuyGoodResponse{}, nil
	}
	payOrderId, err := res.LastInsertId()
	if err != nil {
		log.Printf("在BuyGood中获取创建的支付订单id失败:%s", err)
		return models.BuyGoodResponse{}, nil
	}
	data := models.BuyGoodResponse{
		PayOrderId: payOrderId,
	}
	return data, nil
}

func (q *GoodQuery) CheckPayBeforeNew(ctx context.Context, userId int64) (bool, error) {
	query := "SELECT `id` FROM `pay_orders` WHERE `user_id` = ? AND `status` = ?;"
	row := q.db.QueryRowContext(ctx, query, userId, 0)
	if row.Err() != nil {
		log.Printf("在CheckPayBeforeNew中查询数据库时出错:%s", row.Err().Error())
		return false, row.Err()
	}
	var payOrderId int64
	err := row.Scan(&payOrderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		} else {
			log.Printf("在CheckPayBeforeNew中查询数据库时出错:%s", err)
			return false, err
		}
	}
	return false, nil
}

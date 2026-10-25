package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/utils"
	"github.com/redis/go-redis/v9"
)

type OrderQuery struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewOrderQuery(db *sql.DB, rdb *redis.Client) *OrderQuery {
	return &OrderQuery{
		db:  db,
		rdb: rdb,
	}
}

// 内部工具函数
func (q *OrderQuery) clearStatus(ctx context.Context) {
	key := "status"
	val, _ := q.rdb.SCard(ctx, key).Result()
	if val != 0 {
		q.rdb.Del(ctx, key).Result()
	}
}

// 业务逻辑
func (q *OrderQuery) CheckBeforeNew(ctx context.Context, userId int64) (bool, error) {
	query := "SELECT `id` FROM `orders` WHERE `user_id` = ? AND `status` = ?;"
	row := q.db.QueryRowContext(ctx, query, userId, 0)
	if row.Err() != nil {
		log.Printf("在CheckBeforeNew中查询数据库时出错:%s", row.Err().Error())
		return false, row.Err()
	}
	var orderId int64
	err := row.Scan(&orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		} else {
			log.Printf("在CheckBeforeNew中查询数据库时出错:%s", err)
			return false, err
		}
	}
	return false, nil
}

func (q *OrderQuery) CheckPayBeforeNew(ctx context.Context, userId int64) (bool, error) {
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

func (q *OrderQuery) NewOrder(ctx context.Context, data models.NewOrder) (int64, error) {
	query := "INSERT INTO `orders` (`user_id`, `start_at`, `status`) VALUES ( ?, FROM_UNIXTIME(?), ?);"
	res, err := q.db.ExecContext(ctx, query, data.UserId, data.StartAt, data.Status)
	if err != nil {
		log.Printf("在NewOrder中插入新数据时出现错误:%s", err.Error())
		return 0, err
	}
	orderId, err := res.LastInsertId()
	if err != nil {
		log.Printf("在NewOrder中插入新数据时出现错误:%s", err.Error())
		return 0, err
	}
	// 清除status缓存
	q.clearStatus(ctx)
	return orderId, nil
}

func (q *OrderQuery) GetOrderDataFromCache(ctx context.Context, orderId int64) (models.OrderResponse, error) {
	var data models.OrderResponse
	val, err := q.rdb.Get(ctx, strconv.Itoa(int(orderId))).Result()
	if err != nil {
		return models.OrderResponse{}, nil
	}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		log.Printf("在Redis中查询数据时出错:%s", err)
		return models.OrderResponse{}, nil
	}
	return data, nil
}

func (q *OrderQuery) GetOrderDataById(ctx context.Context, orderId int64) (models.OrderResponse, error) {
	var data models.OrderResponse
	query := "SELECT `id`,`user_id`,`start_at`,`end_at`,`status`,`price` FROM `orders` WHERE `id` = ?;"
	row := q.db.QueryRowContext(ctx, query, orderId)
	if row.Err() != nil {
		log.Printf("在GetOrderDataById中查询数据库时出错:%s", row.Err().Error())
		return models.OrderResponse{}, row.Err()
	}
	// 有可能为空的值单独处理
	var (
		endAt sql.NullTime
		price sql.NullInt64
	)
	err := row.Scan(&data.Id, &data.UserId, &data.StartAt, &endAt, &data.Status, &price)
	if err != nil {
		if err == sql.ErrNoRows {
			// 不存在数据
			return models.OrderResponse{}, models.ErrNoDataExist
		}
		log.Printf("在GetOrderDataById中查询数据库时出错:%s", err.Error())
		return models.OrderResponse{}, err
	}
	// 验证是否存在
	if price.Valid {
		data.Price = price.Int64
	}
	if endAt.Valid {
		data.EndAt = endAt.Time
	}
	// 将所获得到的值存入Redis
	valKey := strconv.Itoa(int(orderId))
	val, err := json.Marshal(data)
	if err != nil {
		log.Printf("在GetOrderDataById中序列化json时出错:%s", err.Error())
		return models.OrderResponse{}, err
	}
	_, err = q.rdb.Set(ctx, valKey, val, 5*time.Minute).Result()
	if err != nil {
		log.Printf("在GetOrderDataById中向Redis插入数据时出错:%s", err.Error())
		return models.OrderResponse{}, nil
	}
	return data, nil
}

func (q *OrderQuery) GetOrderDataByUserId(ctx context.Context, userId int64, status int64, limit int64, offset int64) ([]models.OrderResponse, error) {
	var datas []models.OrderResponse
	var query string
	var rows *sql.Rows
	var err error
	if status < -1 {
		// 当status小于-1时,代表查询不考虑状态
		query = "SELECT `id`, `user_id`, `start_at`, `end_at`, `status`, `price` FROM `orders` WHERE `user_id` = ? ORDER BY update_at DESC LIMIT ? OFFSET ?;"
		// log.Printf("所使用的查询语句为:%s", query)
		rows, err = q.db.QueryContext(ctx, query, userId, limit, offset)
	} else {
		query = "SELECT `id`, `user_id`, `start_at`, `end_at`, `status`, `price` FROM `orders` WHERE `user_id` = ? AND `status` = ? ORDER BY update_at DESC LIMIT ? OFFSET ?;"
		// log.Printf("所使用的查询语句为:%s", query)
		rows, err = q.db.QueryContext(ctx, query, userId, status, limit, offset)
	}
	if err != nil {
		log.Printf("在GetOrderDataByUserId中查询数据库时出错:%s", err.Error())
		return nil, err
	}
	if rows.Err() != nil {
		log.Printf("在GetOrderDataByUserId中查询数据库时出错:%s", rows.Err())
		return nil, rows.Err()
	}
	defer rows.Close()
	for rows.Next() {
		var t models.OrderResponse
		// 对可能为空的值作处理
		var price sql.NullInt64
		var endAt sql.NullTime
		err := rows.Scan(&t.Id, &t.UserId, &t.StartAt, &endAt, &t.Status, &price)
		if err != nil {
			log.Printf("在GetOrderDataByUserId中查询数据库时出错:%s", err.Error())
			return nil, nil
		}
		if price.Valid {
			t.Price = price.Int64
		}
		if endAt.Valid {
			t.EndAt = endAt.Time
		}
		// log.Println(t)
		datas = append(datas, t)
	}
	return datas, nil
}

func (q *OrderQuery) FinishOrder(ctx context.Context, orderId int64) (int64, error) {
	// 计算总价,修改订单状态,生成支付订单
	var startTime time.Time
	var userId int64
	endTime := time.Now()
	status := 1
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("在FinishOrder中开启事务时出错:%s", err.Error())
		return 0, err
	}
	// 具体事务
	// 首先获取订单数据
	query := "SELECT `start_at`, `user_id` FROM `orders` WHERE `id` = ?;"
	row := tx.QueryRowContext(ctx, query, orderId)
	if row.Err() != nil {
		tx.Rollback()
		return 0, row.Err()
	}
	err = row.Scan(&startTime, &userId)
	if err != nil {
		return 0, err
	}
	totalPrice, err := utils.CalcTotalPrice(startTime, endTime)
	if err != nil {
		return 0, err
	}
	// 更新订单数据
	query = "UPDATE `orders` SET `end_at` = ?,`status` = ?,`price` = ? WHERE id = ?;"
	res, err := tx.ExecContext(ctx, query, endTime, status, totalPrice, orderId)
	if err != nil {
		tx.Rollback()
		log.Printf("在FinishOrder中更新数据时出现错误:%s", err.Error())
		return 0, nil
	}
	// 涉及更新订单,需要处理缓存中的数据
	valKey := strconv.Itoa(int(orderId))
	q.rdb.Del(ctx, valKey)
	// 创建支付订单
	query = "INSERT INTO pay_orders (`user_id`, `type`, `order_no`, `status`, `description`, `total_fee`) VALUES (?, ?, ?, ?, ?, ?);"
	res, err = tx.ExecContext(ctx, query, userId, "order", orderId, 0, "签到结算", totalPrice)
	if err != nil {
		log.Printf("在FinishOrder中创建支付订单时出错:%s", err)
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		log.Printf("在FinishOrder中提交事务时出错:%s", err.Error())
		return 0, err
	}
	// 清除status缓存
	q.clearStatus(ctx)
	payOrderId, err := res.LastInsertId()
	if err != nil {
		log.Printf("在FinishOrder中提交事务时出错:%s", err)
		return 0, err
	}
	if payOrderId > 0 {
		return payOrderId, nil
	} else {
		return 0, nil
	}
}

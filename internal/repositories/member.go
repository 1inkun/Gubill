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

type MemberQuery struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewMemberQuery(db *sql.DB, rdb *redis.Client) *MemberQuery {
	return &MemberQuery{
		db:  db,
		rdb: rdb,
	}
}

func (q *MemberQuery) GetMemberData(ctx context.Context, userId int64) (models.MemberData, error) {
	var data models.MemberData
	query := "SELECT `id`, `user_id`, `start_at`, `end_at`, `status` FROM members WHERE `user_id` = ?;"
	row := q.db.QueryRowContext(ctx, query, userId)
	if row.Err() != nil {
		log.Printf("在GetMemberData中查询数据库时出错:%s", row.Err())
		return models.MemberData{}, row.Err()
	}
	err := row.Scan(&data.Id, &data.UserId, &data.StartAt, &data.EndAt, &data.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.MemberData{}, nil
		}
		log.Printf("在GetMemberData中查询数据库时出错:%s", err)
		return models.MemberData{}, err
	}
	return data, nil
}

func (q *MemberQuery) GetMemberDataFromCache(ctx context.Context, key string) (models.MemberData, error) {
	val, err := q.rdb.Get(ctx, key).Result()
	if err != nil {
		log.Printf("在GetMemberData中查询缓存时出现错误:%s", err)
		return models.MemberData{}, err
	}
	// 反序列化
	var data models.MemberData
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		log.Printf("在GetMemberData中查询缓存时出现错误:%s", err)
		return models.MemberData{}, err
	}
	return data, nil
}

func (q *MemberQuery) SetCache(ctx context.Context, key string, value string) error {
	_, err := q.rdb.Set(ctx, key, value, 10*time.Minute).Result()
	if err != nil {
		return err
	}
	return nil
}

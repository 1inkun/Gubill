package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/1inkun/billApp/internal/models"
	"github.com/redis/go-redis/v9"
)

type StatusQuery struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewStatusQuery(db *sql.DB, rdb *redis.Client) *StatusQuery {
	return &StatusQuery{
		db:  db,
		rdb: rdb,
	}
}

func (q *StatusQuery) GetCurrentUser(ctx context.Context) ([]int64, error) {
	var data []int64
	query := "SELECT user_id FROM orders WHERE status = 0;"
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("在GetCurrentUser中查询数据库时出错:%s", err.Error())
		return nil, err
	}
	if rows.Err() != nil {
		log.Printf("在GetCurrentUser中查询数据库时出错:%s", rows.Err())
		return nil, rows.Err()
	}
	defer rows.Close()
	for rows.Next() {
		var i int64
		if err := rows.Scan(&i); err != nil {
			log.Printf("在GetCurrentUser中查询数据库时出错:%s", err.Error())
			return nil, err
		}
		data = append(data, i)
	}
	return data, nil
}

func (q *StatusQuery) GetUserData(ctx context.Context, userId int64) (models.CurrentUser, error) {
	var data models.CurrentUser
	query := "SELECT `username`, `nickname`, `avatar_file` FROM `users` WHERE `id` = ?;"
	row := q.db.QueryRowContext(ctx, query, userId)
	if row.Err() != nil {
		log.Printf("在GetUserData中查询数据库时出错:%s", row.Err())
		return models.CurrentUser{}, row.Err()
	}
	// 对可能为空的值单独处理
	var (
		nickname   sql.NullString
		avatarFile sql.NullString
	)
	err := row.Scan(&data.Username, &nickname, &avatarFile)
	if err != nil {
		log.Printf("在GetUserData中查询数据库时出错:%s", err.Error())
		return models.CurrentUser{}, err
	}
	if nickname.Valid {
		data.Nickname = nickname.String
	}
	if avatarFile.Valid {
		data.AvatarFile = avatarFile.String
	}
	// 将查询到的值添加到缓存中
	val, err := json.Marshal(data)
	if err != nil {
		log.Printf("Status缓存创建失败:%s", err)
		return data, nil
	}
	q.rdb.SAdd(ctx, "status", val).Result()
	return data, nil
}

func (q *StatusQuery) GetUserDataFromCache(ctx context.Context) ([]models.CurrentUser, error) {
	var datas []models.CurrentUser
	key := "status"
	len, err := q.rdb.SCard(ctx, key).Result()
	if len == 0 {
		// 集合长度为零表示没有缓存的数据
		return nil, nil
	}
	vals, err := q.rdb.SMembers(ctx, key).Result()
	if err != nil {
		log.Printf("在GetUserDataFromCache中查询Redis时遇到错误:%s", err)
		return nil, err
	}
	for _, v := range vals {
		var i models.CurrentUser
		err := json.Unmarshal([]byte(v), &i)
		if err != nil {
			log.Printf("在GetUserDataFromCache中序列化json时遇到错误:%s", err)
			return nil, err
		}
		datas = append(datas, i)
	}
	return datas, err
}

package repositories

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/1inkun/billApp/internal/models"
	"github.com/redis/go-redis/v9"
)

type UserQuery struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewUserQuery(db *sql.DB, rdb *redis.Client) *UserQuery {
	return &UserQuery{
		db:  db,
		rdb: rdb,
	}
}

func (q *UserQuery) Register(ctx context.Context, data models.RegisterData) error {
	query := "INSERT INTO `users` (`username`, `nickname`, `password`, `email`, `role`) VALUES (?, ?, ?, ?,'user');"
	_, err := q.db.ExecContext(ctx, query, data.Username, data.Nickname, data.PasswdHash, data.Email)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (q *UserQuery) GetUserPasswdHashByUername(ctx context.Context, username string) (string, error) {
	var passwdHash string
	query := "SELECT `password` FROM `users` WHERE `username` = ?;"
	row := q.db.QueryRowContext(ctx, query, username)
	if row.Err() != nil {
		log.Println(row.Err())
		return "", row.Err()
	}
	err := row.Scan(&passwdHash)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return passwdHash, nil
}

func (q *UserQuery) GetClaimsData(ctx context.Context, username string) (models.JwtClaims, error) {
	var claims models.JwtClaims
	claims.Username = username
	query := "SELECT `id`,`role` FROM `users` WHERE `username` = ?;"
	row := q.db.QueryRow(query, username)
	if row.Err() != nil {
		log.Println(row.Err())
		return models.JwtClaims{}, row.Err()
	}
	row.Scan(&claims.UserId, &claims.Role)
	return claims, nil
}

func (q *UserQuery) AddTokenToRedis(ctx context.Context, username string, tokenStr string) error {
	// 先确定redis中是否已存有token
	val, _ := q.rdb.Get(ctx, username).Result()
	// log.Println(val)
	if val != "" {
		_, err := q.rdb.Del(ctx, username).Result()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	// 设置token
	exphours, err := strconv.Atoi(os.Getenv("exptime"))
	if err != nil {
		log.Printf("读取配置文件时出现错误:%s", err.Error())
		return err
	}
	exp := time.Hour * time.Duration(exphours)
	_, err = q.rdb.Set(ctx, username, tokenStr, exp).Result()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

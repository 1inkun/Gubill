package repository

import (
	"os"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabaseConnect() (*gorm.DB, error) {
	sqliteUrl := os.Getenv("SQLiteUrl")
	db, err := gorm.Open(sqlite.Open(sqliteUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

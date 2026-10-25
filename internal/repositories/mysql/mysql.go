package mysql

import (
	"database/sql"
	"log"
	"os"

	"github.com/1inkun/billApp/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewDBConection(cfg config.MySqlConfig) (*sql.DB, error) {
	var dsn = fmtdsn(cfg)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("出现致命错误:%s", err.Error())
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接超时:%s", err.Error())
	}
	if err = InitDatabase(db); err != nil {
		return nil, err
	}
	// 设置数据库连接池
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	// 数据库连接成功
	log.Println("MySql数据库连接成功!")
	return db, nil
}

func fmtdsn(cfg config.MySqlConfig) string {
	var dsn string
	dsn = cfg.User + ":" + cfg.Password + "@" + "tcp(" + cfg.Host + ":" + cfg.Port + ")" + "/" + cfg.Database + "?multiStatements=true&parseTime=true"
	// log.Println(dsn)
	return dsn
}

func InitDatabase(db *sql.DB) error {
	// 使用schema.sql来初始化数据库
	file, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatalf("读取schema文件失败:%s", err.Error())
		return err
	}
	query := string(file)
	// log.Println(query)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("初始化数据库失败:%s", err.Error())
		return err
	}
	return nil
}

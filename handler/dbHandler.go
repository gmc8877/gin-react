package handler

import (
	"database/sql"
	"log"
	"time"
	_ "github.com/go-sql-driver/mysql"
)
var DB *sql.DB

func InitDB() {
	DB, _ = sql.Open("mysql", "root:root@tcp(localhost:3306)/demo01?charset=utf8&parseTime=True&loc=Local") // 使用本地时间，即东八区，北京时间
	// set pool params
	DB.SetMaxOpenConns(2000)
	DB.SetMaxIdleConns(1000)
	DB.SetConnMaxLifetime(time.Minute * 60) // mysql default conn timeout=8h, should < mysql_timeout
	err := DB.Ping()
	if err != nil {
		log.Println("database init failed, err: ", err)
	}
	log.Println("mysql conn pool has initiated.")
}
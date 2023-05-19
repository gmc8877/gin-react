package handler

import (
	"DEMO01/tools"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
)



var DB *sql.DB

func InitDB() {
	var mysqlconfig = new(Config)
	err := ini.MapTo(mysqlconfig, "./conf/config.ini")
	tools.CheckErr(err)
	sql_l:= fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local", mysqlconfig.User,mysqlconfig.Mysqlconfig.Password, mysqlconfig.Host, mysqlconfig.Port,mysqlconfig.Dbname )
	DB, _ = sql.Open(mysqlconfig.Db, sql_l) // 使用本地时间，即东八区，北京时间
	// set pool params
	DB.SetMaxOpenConns(2000)
	DB.SetMaxIdleConns(1000)
	DB.SetConnMaxLifetime(time.Minute * 60) // mysql default conn timeout=8h, should < mysql_timeout
	err = DB.Ping()
	if err != nil {
		log.Println("database init failed, err: ", err)
	}
	log.Println("mysql conn pool has initiated.")
}
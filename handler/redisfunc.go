package handler

import (
	"DEMO01/tools"
	"context"
	"fmt"

	"github.com/go-ini/ini"
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func RedisInit() {

	var ctx = context.Background()

	//cionfig.ini文件进行配置
	var redisconfig = new(Config)
	err := ini.MapTo(redisconfig, "./conf/config.ini")
	tools.CheckErr(err)
	//初始化redis，连接地址和端口，密码，数据库名称
	Rdb = redis.NewClient(&redis.Options{
		Addr:         redisconfig.Address,
		Password:     "1",
		DB:           0,
		MinIdleConns: 1,
		PoolSize:     1000,
	})
	res, err := Rdb.Ping(ctx).Result()
	fmt.Println(res)
	tools.CheckErr(err)
}

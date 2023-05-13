package handler
import "github.com/go-redis/redis/v8"

var Rdb *redis.Client
func RedisInit() {
	//初始化redis，连接地址和端口，密码，数据库名称
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "1",
		DB:       0,
		MinIdleConns: 1,
    PoolSize:     1000,
	})
}

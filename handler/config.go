package handler


type Config struct {
	Mysqlconfig `ini:"mysql"`
	Redisconfig `ini:"redis"`
}

type Mysqlconfig struct {
	Db string `ini:"db"`
	Host string `ini:"host"`
	Port int `ini:"port"`
	User string `ini:"user"`
	Password string `ini:"password"`
	Dbname string `ini:"dbname"`
}

type Redisconfig struct {
	Address string `ini:"address"`
	Password string `ini:"password"`
}
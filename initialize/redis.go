package initialize

import (
	"chatgpt-go/db"
	"chatgpt-go/global"
	"github.com/go-redis/redis"
)

func InitRedis() {
	db.RedisDb = redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Addr,
		Password: global.Config.Redis.Password,
		DB:       0,
	})
}

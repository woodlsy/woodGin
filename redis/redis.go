package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/woodlsy/woodGin/config"
	"github.com/woodlsy/woodGin/helper"
	"github.com/woodlsy/woodGin/log"
)

var Redis *redis.Client

func Enabled() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Configs.Redis.Host, config.Configs.Redis.Port),
		Password: config.Configs.Redis.Password, // no password set
		DB:       config.Configs.Redis.Db,       // use default DB
	})

	if err := Redis.Ping().Err(); err != nil { //心跳测试
		log.Logger.Error("redis连接失败", err, config.Configs.Redis)
		errMsg := "failed to init redis"
		panic(errMsg)
	}
	fmt.Println(helper.Now(), "redis 连接成功")
}

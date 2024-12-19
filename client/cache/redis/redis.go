package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/woodlsy/woodGin/client/cache"
	"github.com/woodlsy/woodGin/config"
	"github.com/woodlsy/woodGin/helper"
	"github.com/woodlsy/woodGin/log"
	"time"
)

type Cache struct {
	redis  *redis.Client
	prefix string
}

func init() {
	cache.RegisterCache("redis", NewRedisCache)
}

func (c *Cache) Enabled() error {
	c.redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Configs.Redis.Host, config.Configs.Redis.Port),
		Password: config.Configs.Redis.Password, // no password set
		DB:       config.Configs.Redis.Db,       // use default DB
	})

	if err := c.redis.Ping().Err(); err != nil { //心跳测试
		log.Logger.Error("redis连接失败", err, config.Configs.Redis)
		errMsg := "failed to init redis"
		return errors.New(errMsg)
	}
	fmt.Println(helper.Now(), "redis 连接成功")
	return nil
}

func (c *Cache) SetPrefix() {
	c.prefix = config.Configs.Redis.Prefix
}

func NewRedisCache() cache.Cache {
	return &Cache{}
}

func (c *Cache) Exists(key string) bool {
	defer c.Close()
	if c.prefix != "" {
		key = helper.Join("", c.prefix, key)
	}
	v, err := c.redis.Exists(key).Result()
	if err != nil {
		log.Logger.Error("【redis】【Exists】key:", key, "error:", err)
		return false
	}
	return v > 0
}

func (c *Cache) Get(key string) string {
	defer c.Close()
	if c.prefix != "" {
		key = helper.Join("", c.prefix, key)
	}
	value, err := c.redis.Get(key).Result()
	if err != nil {
		log.Logger.Error("【redis】【Get】key:", key, "value:", value, "error:", err)
		return ""
	}
	return value
}

//
//func Set(key string, value interface{}, ttl int) bool {
//	err := woodlsy.Redis.Set(key, value, time.Second*time.Duration(ttl)).Err()
//	if err != nil {
//		log.Logger.Error("【redis】【SetEx】key:", key, "value:", value, "ttl:", ttl, "error:", err)
//		return false
//	}
//	return true
//}

func (c *Cache) SetEx(key string, ttl int, value interface{}) bool {
	defer c.Close()
	if c.prefix != "" {
		key = helper.Join("", c.prefix, key)
	}
	err := c.redis.Set(key, value, time.Second*time.Duration(ttl)).Err()
	if err != nil {
		log.Logger.Error("【redis】【SetEx】key:", key, "value:", value, "ttl:", ttl, "error:", err)
		return false
	}
	return true
}

func (c *Cache) Del(key string) bool {
	defer c.Close()
	if c.prefix != "" {
		key = helper.Join("", c.prefix, key)
	}
	err := c.redis.Del(key).Err()
	if err != nil {
		log.Logger.Error("【redis】【Del】key:", key, "error:", err)
		return false
	}
	return true
}

func (c *Cache) Close() error {
	//return c.redis.Close()
	return nil
}

func (c *Cache) Expire(key string, ttl int) bool {
	defer c.Close()
	if c.prefix != "" {
		key = helper.Join("", c.prefix, key)
	}
	err := c.redis.Expire(key, time.Second*time.Duration(ttl)).Err()
	if err != nil {
		log.Logger.Error("【redis】【Expire】key:", key, "ttl:", ttl, "error:", err)
		return false
	}
	return true
}

func (c *Cache) Incr(key string) int64 {
	if c.prefix != "" {
		key = helper.Join("", c.prefix, key)
	}
	value, err := c.redis.Incr(key).Result()
	if err != nil {
		log.Logger.Error("【redis】【Incr】key:", key, "error:", err)
		return 0
	}
	return value
}

func (c *Cache) Decr(key string) int64 {
	if c.prefix != "" {
		key = helper.Join("", c.prefix, key)
	}
	value, err := c.redis.Decr(key).Result()
	if err != nil {
		log.Logger.Error("【redis】【decr】key:", key, "error:", err)
		return 0
	}
	return value
}

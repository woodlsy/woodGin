package cache

import "errors"

type Cache interface {
	Enabled() error
	Exists(key string) bool
	Get(key string) string
	SetEx(key string, ttl int, value interface{}) bool
	Del(key string) bool
	Incr(key string) int64
	Decr(key string) int64
	Expire(key string, ttl int) bool
	Close() error
	SetPrefix()
}

type Instance func() Cache

var adapters = make(map[string]Instance)

func RegisterCache(name string, adapter Instance) {
	if _, ok := adapters[name]; !ok {
		adapters[name] = adapter
	}
}

func NewCache(name string) (adapter Cache, err error) {
	adapterFunc, ok := adapters[name]
	if !ok {
		err = errors.New("不支持的缓存")
		return
	}
	adapter = adapterFunc()
	err = adapter.Enabled()
	if err != nil {
		adapter = nil
	}
	return
}

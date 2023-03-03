package db

import "errors"

type Orm interface {
	Enabled()
	SetDatabase(name string) *Orm
	GetAll(m interface{}, where map[string]interface{}, orderBy string, fields string) (int64, error)
}

type Instance func() Orm

var adapters = make(map[string]Instance)

func RegisterDb(name string, adapter Instance) {
	if _, ok := adapters[name]; !ok {
		adapters[name] = adapter
	}
}

func NewDb(name string) (adapter Orm, err error) {
	adapterFunc, ok := adapters[name]
	if !ok {
		err = errors.New("不支持的数据库类型")
		return
	}
	adapter = adapterFunc()
	adapter.Enabled()

	return
}

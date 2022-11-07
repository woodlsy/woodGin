package db

import (
	"github.com/woodlsy/woodGin/log"
	"gorm.io/gorm"
)

type Orm struct {
	DbName string
}

type OrmInterface interface {
	GetOne(interface{}, map[string]interface{}, string)
	GetList(interface{}, map[string]interface{}, string, int, int)
	GetAll(interface{}, map[string]interface{}, string)
	Count()
}

var db map[string]*gorm.DB

func OrmInit() {
	db = connect()
}

func (o Orm) conn() *gorm.DB {
	if _, ok := db[o.DbName]; !ok {
		log.Logger.Errorf("切换的数据库%s不存在", o.DbName)
		panic("切换的数据库不存在")
	}
	return db[o.DbName]
}

func (o Orm) Model(value interface{}) *gorm.DB {
	return o.conn().Model(value)
}

func (o Orm) Insert(value interface{}) *gorm.DB {
	return o.conn().Create(value)
}

func (o Orm) Update(value interface{}, data map[string]interface{}, where map[string]interface{}) int64 {
	result := o.conn().Model(value).Where(where).Updates(data)
	return result.RowsAffected
}
func (o Orm) Deleted(value interface{}, where map[string]interface{}) int64 {
	result := o.conn().Where(where).Delete(value)
	return result.RowsAffected
}

func (o Orm) Where(query interface{}, args ...interface{}) *gorm.DB {
	return o.conn().Where(query, args...)
}

func (o Orm) Select(query interface{}, args ...interface{}) *gorm.DB {
	return o.conn().Select(query, args...)
}

func (o Orm) Source() *gorm.DB {
	return o.conn()
}

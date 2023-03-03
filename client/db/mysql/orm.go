package mysql

import (
	"github.com/woodlsy/woodGin/helper"
	"github.com/woodlsy/woodGin/log"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type Orm struct {
	DatabaseName string
	conn         *gorm.DB
}

type OrmInterface interface {
	GetOne(interface{}, map[string]interface{}, string)
	GetList(interface{}, map[string]interface{}, string, int, int)
	GetAll(interface{}, map[string]interface{}, string)
	Count()
}

var db map[string]*gorm.DB

func Enabled() {
	db = connect()
}

func (o *Orm) ChangeConn() *Orm {
	if o.conn != db[o.DatabaseName] {
		if _, ok := db[o.DatabaseName]; !ok {
			//log.Logger.Errorf("切换的数据库%s不存在", name)
			helper.P(db)
			panic("切换的数据库不存在")
		}
		o.conn = db[o.DatabaseName]
	}
	return o
}

//
// Updates
// @Description: 根据主键ID更新
// @receiver o
// @param value
// @param updateFields
// @return int64
//
func (o Orm) Updates(value interface{}, updateFields []string) int64 {
	orm := o.Source().Model(value)
	if len(updateFields) > 0 {
		orm.Select(updateFields)
	}
	result := orm.Updates(value)
	return result.RowsAffected
}

//
// UpdatesWhere
// @Description: 根据条件更新
// @receiver o
// @param value
// @param data
// @param where
// @return int64
//
func (o Orm) UpdatesWhere(value interface{}, data map[string]interface{}, where map[string]interface{}) int64 {
	result := o.SqlCondition(value, where, "", 0, 0, "").Updates(data)
	return result.RowsAffected
}

func (o Orm) Deleted(value interface{}, where map[string]interface{}) int64 {
	result := o.SqlCondition(value, where, "", 0, 0, "").Delete(value)
	return result.RowsAffected
}

func (o Orm) Insert(value interface{}) {
	result := o.Source().Create(value)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Logger.Error("新增记录失败", value, result.Error)
	}
}

func (o Orm) Source() *gorm.DB {
	if o.conn == nil {
		o.ChangeConn()
	}
	return o.conn
}

func (o Orm) Transaction() *gorm.DB {
	return o.Source().Begin()
}

//
// GetOne
// @Description: 获取单条记录
// @param m
// @param where
// @param orderBy
// @return int64
// @return error
//
//func (o Orm) GetOne(m interface{}, where map[string]interface{}, orderBy string, fields string) (int64, error) {
//	o := o.sqlCondition(m, where, orderBy, 0, 0, fields)
//	result := o.Order(orderBy).Find(&m)
//	return result.RowsAffected, result.Error
//}

//
// sqlCondition
// @Description: sql 条件组装
// @param m
// @param where
// @param orderBy
// @param offset
// @param limit
// @param fields
// @return *gorm.DB
//
func (o Orm) SqlCondition(m interface{}, where map[string]interface{}, orderBy string, offset int, limit int, fields string) *gorm.DB {
	tx := o.Source().Model(m)

	if len(where) > 0 {
		whereSqlString, whereValueArray := parseWhere(where)
		if whereSqlString != "" {
			tx.Where(whereSqlString, whereValueArray...)
		}
	}
	if orderBy != "" {
		tx.Order(orderBy)
	}
	if limit != 0 {
		tx.Offset(offset).Limit(limit)
	}
	if fields != "" {
		tx.Select(fields)
	}
	return tx
}

//
// parseWhere
// @Description:解析where条件
//	map[string]interface{}{
//	"id":        2,                                                    // id =2
//	"cid":       []interface{}{"in", []int{1, 3}},                                         // cid in (1,2)
//	"pid":       []interface{}{"!=", 3},                               // pid != 3
//	"name":      []interface{}{"like", "%%大三%%"},                      // name like '%大三%'
//	"create_at": []interface{}{"between", "2022-08-01", "2022-08-02"}, // create_at between '2022-08-01' and '2022-08-02'
//	"bid":       []interface{}{"not in", []int{1, 2}},                 // bid not in (1, 2)
//	"or": map[string]interface{}{ // or链接
//	"id":        1,                                                    // id =2
//	"cid":       []interface{}{"in", []int{1, 3}},                                         // cid in (1,2)
//	"pid":       []interface{}{"!=", 4},                               // pid != 3
//	"name":      []interface{}{"like", "%%大三2%%"},                     // name like '%大三%'
//	"create_at": []interface{}{"between", "2022-08-03", "2022-08-04"}, // create_at between '2022-08-01' and '2022-08-02'
//	"bid":       []interface{}{"not in", []int{1, 3}},                 // bid not in (1, 2)
//	},
//	}
// @param where
// @return string
// @return []interface{}
//
func parseWhere(where map[string]interface{}) (string, []interface{}) {
	var whereSql []string
	var whereValue []interface{}
	var childWhereValue []interface{}
	var orWhereSqlString string
	for key, value := range where {
		if key == "or" && reflect.TypeOf(value).Kind() == reflect.Map {
			orWhereSqlString, childWhereValue = parseWhere(value.(map[string]interface{}))
		} else {
			switch vv := value.(type) {
			case []interface{}:
				switch strings.ToLower(vv[0].(string)) {
				case "!=", ">", ">=", "<", "<=", "like", "in", "not in":
					whereSql = append(whereSql, helper.Join(" ", key, vv[0].(string), "?"))
					whereValue = append(whereValue, vv[1])
				case "between":
					whereSql = append(whereSql, helper.Join(" ", key, vv[0].(string), "?", "and", "?"))
					whereValue = append(whereValue, vv[1], vv[2])
				}
			default:
				whereSql = append(whereSql, helper.Join(" ", key, "=", "?"))

				whereValue = append(whereValue, value)
			}
		}
	}
	whereSqlString := strings.Join(whereSql, " AND ")

	var tmpSqlArray []string
	if len(whereSqlString) > 0 {
		tmpSqlArray = append(tmpSqlArray, helper.Join(" ", "(", whereSqlString, ")"))
	}
	if len(orWhereSqlString) > 0 {
		tmpSqlArray = append(tmpSqlArray, helper.Join(" ", "(", orWhereSqlString, ")"))
		return strings.Join(tmpSqlArray, " OR "), append(whereValue, childWhereValue...)
	}
	return whereSqlString, whereValue
}

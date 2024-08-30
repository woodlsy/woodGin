package demo

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/woodlsy/woodGin/client/db/mysql"
	"github.com/woodlsy/woodGin/helper"
	"gorm.io/gorm"
	"time"
)

var cst *time.Location

func init() {
	var err error
	cst, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("load Asia/Shanghai err:", err)
	}

}

type LocalDate time.Time

const EmptyLocalDate = "0001-01-01"

func (t *LocalDate) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	if &t == nil || t.IsZero() || t.String() == "1990-01-01" || t.String() == "0001-01-01" {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02"))), nil
}

func (t LocalDate) String() string {
	return time.Time(t).In(cst).Format("2006-01-02")
}

func (t LocalDate) IsZero() bool {
	tTime := time.Time(t)
	return tTime.IsZero() || tTime.String() == "0001-01-01"
}

func (t LocalDate) Value() (driver.Value, error) {
	return []byte(time.Time(t).Format("2006-01-02")), nil
}

func (t *LocalDate) Scan(value interface{}) error {
	if value == nil {
		*t = LocalDate(time.Time{})
		return nil
	}

	var err error
	switch v := value.(type) {
	case time.Time:
		*t = LocalDate(v.In(cst))
	case []byte:
		parsedTime, err := time.ParseInLocation("2006-01-02", string(v), cst)
		if err != nil {
			return err
		}
		*t = LocalDate(parsedTime)
	case string:
		parsedTime, err := time.ParseInLocation("2006-01-02", v, cst)
		if err != nil {
			return err
		}
		*t = LocalDate(parsedTime)
	default:
		return fmt.Errorf("cannot scan type %T into LocalDate", value)
	}
	return err
}

func NewLocalDate(dateStr string) LocalDate {
	defaultDate := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	localDate := LocalDate(defaultDate)
	if dateStr == "" {
		return localDate
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return localDate
	}
	localDate = LocalDate(t)
	return localDate
}

type LocalDateTime time.Time

const EmptyLocalDateTime = "0001-01-01 00:00:00"

func (t *LocalDateTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	if &t == nil || t.IsZero() || t.String() == "1990-01-01 00:00:00" || t.String() == "0001-01-01 00:00:00" {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (t LocalDateTime) IsZero() bool {
	tTime := time.Time(t)
	return tTime.IsZero() || tTime.String() == "0001-01-01 00:00:00"
}

func (t LocalDateTime) Value() (driver.Value, error) {
	return []byte(time.Time(t).Format("2006-01-02 15:04:05")), nil
}

func (t LocalDateTime) String() string {
	return time.Time(t).In(cst).Format("2006-01-02 15:04:05")
}

func (t *LocalDateTime) Scan(value interface{}) error {

	if value == nil {
		*t = LocalDateTime(time.Time{})
		return nil
	}

	var err error
	switch v := value.(type) {
	case time.Time:
		*t = LocalDateTime(v.In(cst))
	case []byte:
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", string(v), cst)
		if err != nil {
			return err
		}
		*t = LocalDateTime(parsedTime)
	case string:
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", v, cst)
		if err != nil {
			return err
		}
		*t = LocalDateTime(parsedTime)
	default:
		return fmt.Errorf("cannot scan type %T into LocalDateTim", value)
	}
	return err
}

func NewLocalDateTime(dateStr string) LocalDateTime {
	defaultDate := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	localDateTime := LocalDateTime(defaultDate)
	if dateStr == "" {
		return localDateTime
	}

	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return localDateTime
	}
	localDateTime = LocalDateTime(t)
	return localDateTime
}

type Model struct {
	Id       int64         `gorm:"primarykey" json:"id,omitempty"`
	CreateAt LocalDateTime `json:"createAt,omitempty"`
	UpdateAt LocalDateTime `json:"updateAt,omitempty"`
	CreateBy int64         `json:"createBy,omitempty"`
	UpdateBy int64         `json:"updateBy,omitempty"`
}

/************模板 start**************/

var Orm mysql.Orm

type ModelTemp struct {
	Model
}

func (m *ModelTemp) GetAll(where map[string]interface{}, orderBy string, fields string) []*ModelTemp {
	return m.GetList(where, orderBy, 0, 0, fields)
}

func (m *ModelTemp) GetList(where map[string]interface{}, orderBy string, offset int, limit int, fields string) []*ModelTemp {
	result := make([]*ModelTemp, 0)
	Orm.SqlCondition(m, where, orderBy, offset, limit, fields).Find(&result)
	return result
}

func (m *ModelTemp) GetOne(where map[string]interface{}, orderBy string, fields string) (row *ModelTemp) {
	result := Orm.SqlCondition(m, where, orderBy, 0, 1, fields).Find(&row)
	if result.RowsAffected == 0 {
		return nil
	}
	return row
}

func (m *ModelTemp) GetById(id int64, fields string) *ModelTemp {
	if id == 0 {
		return nil
	}
	return m.GetOne(map[string]interface{}{"id": id}, "", fields)
}

func (m *ModelTemp) GetCount(where map[string]interface{}) (row int64) {
	Orm.SqlCondition(m, where, "", 0, 0, "").Count(&row)
	return
}

func (m *ModelTemp) GetSum(where map[string]interface{}, field string) float64 {
	var total sql.NullFloat64
	Orm.SqlCondition(m, where, "", 0, 0, fmt.Sprintf("COALESCE(SUM(%s), 0) as total", field)).Scan(&total)
	if total.Valid {
		return total.Float64
	} else {
		return 0
	}
}

func (m *ModelTemp) CustomReturnOne(result interface{}, where map[string]interface{}, orderBy string, fields string) interface{} {
	Orm.SqlCondition(m, where, orderBy, 0, 1, fields).First(result)
	return result
}

func (m *ModelTemp) CustomReturnAll(result interface{}, where map[string]interface{}, orderBy string, fields string) interface{} {
	m.CustomReturnList(result, where, orderBy, 0, 0, fields)
	return result
}

func (m *ModelTemp) CustomReturnList(result interface{}, where map[string]interface{}, orderBy string, offset int, limit int, fields string) interface{} {
	Orm.SqlCondition(m, where, orderBy, offset, limit, fields).Find(result)
	return result
}

func (m *ModelTemp) Insert() int64 {
	if m.CreateAt.IsZero() {
		m.CreateAt = LocalDateTime(time.Now())
		m.UpdateAt = m.CreateAt
	}
	Orm.Insert(&m)
	return m.Id
}

func (m *ModelTemp) Updates(fields []string) int64 {
	if m.UpdateAt.IsZero() {
		m.UpdateAt = LocalDateTime(time.Now())
	}
	return Orm.Updates(&m, fields)
}

func (m *ModelTemp) UpdatesWhere(data map[string]interface{}, where map[string]interface{}) int64 {
	if _, ok := data["update_at"]; !ok {
		data["update_at"] = helper.Now()
	}
	return Orm.UpdatesWhere(&m, data, where)
}

func (m *ModelTemp) Delete() int64 {
	return Orm.Delete(&m)
}

func (m *ModelTemp) DeleteWhere(where map[string]interface{}) int64 {
	return Orm.DeleteWhere(m, where)
}

func (m *ModelTemp) TransactionDelete(db *gorm.DB, where map[string]interface{}) (int64, error) {
	result := Orm.ParseWhere(db, where).Delete(&ModelTemp{})
	return result.RowsAffected, result.Error
}

func (m *ModelTemp) TransactionUpdates(db *gorm.DB, data map[string]interface{}, where map[string]interface{}) (int64, error) {
	result := Orm.ParseWhere(db, where).Model(&m).Updates(data)
	return result.RowsAffected, result.Error
}

// TransactionBatchInsert
//
//	@Description: 事务式批量插入
//	@receiver m
//	@param db
//	@param insertData
//	@return int64
//	@return error
func (m *ModelTemp) TransactionBatchInsert(db *gorm.DB, insertData []*ModelTemp) (int64, error) {
	for _, value := range insertData {
		if value.CreateAt.IsZero() {
			value.CreateAt = LocalDateTime(time.Now())
			value.UpdateAt = value.CreateAt
		}
	}
	result := db.Create(insertData)
	return result.RowsAffected, result.Error
}

// TransactionSave
// @Description: 事务式保存
// @receiver m
// @param db
// @param fields
// @return int
func (m *ModelTemp) TransactionSave(db *gorm.DB, fields ...interface{}) (int64, error) {
	if m.UpdateAt.IsZero() {
		m.UpdateAt = LocalDateTime(time.Now())
	}
	if m.Id == 0 {
		if m.CreateAt.IsZero() {
			m.CreateAt = LocalDateTime(time.Now())
		}
		result := db.Create(&m)
		return m.Id, result.Error
	} else {
		if len(fields) > 0 {
			db = db.Select(fields[0])
		}
		result := db.Updates(&m)
		return result.RowsAffected, result.Error
	}
}

/************模板 end**************/

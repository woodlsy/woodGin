package helper

import (
	"errors"
	"reflect"
)

//
// GetStructField
// @Description: 获取结构体中某个字段的值
// @param input
// @param key
// @return value
// @return err
//
func GetStructField(input interface{}, key string) (value interface{}, err error) {
	rv := reflect.ValueOf(input)
	rt := reflect.TypeOf(input)
	if rt.Kind() != reflect.Struct {
		return value, errors.New("input must be struct")
	}
	val := rv.FieldByName(key)
	emptyValue := reflect.Value{}
	if val == emptyValue {
		return "", errors.New("key no exists")
	}
	return val.Interface(), nil
}

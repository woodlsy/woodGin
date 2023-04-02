package helper

import (
	"errors"
	"reflect"
	"fmt"
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

// extractFieldAsMap 函数，接受一个切片结构体和一个字段名作为参数，返回一个以该字段值为下标的 map
func ExtractFieldAsMap(slice interface{}, fieldName string) map[interface{}]interface{} {
	// 使用反射获取切片的类型和值
	sliceValue := reflect.ValueOf(slice)
	sliceType := sliceValue.Type()

	// 确保传入的参数是一个切片结构体
	if sliceType.Kind() != reflect.Slice || sliceType.Elem().Kind() != reflect.Struct {
		panic("extractFieldAsMap: not a slice of structs")
	}

	// 确保切片中包含指定字段
	field, found := sliceType.Elem().FieldByName(fieldName)
	if !found {
		panic(fmt.Sprintf("extractFieldAsMap: field %s not found", fieldName))
	}

	// 定义一个空的 map，用于存储提取出的字段值
	result := make(map[interface{}]interface{})

	// 遍历切片，提取出指定字段的值作为下标，并将整个结构体作为值存储到 map 中
	for i := 0; i < sliceValue.Len(); i++ {
		structValue := sliceValue.Index(i)
		key := structValue.FieldByIndex(field.Index).Interface()
		result[key] = structValue.Interface()
	}

	return result
}
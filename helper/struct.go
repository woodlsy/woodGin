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
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rv.Type()
	}
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

/**
 * ExtractFieldAsMap
 * @Description: 接受一个切片结构体和一个字段名作为参数，返回一个以该字段值为下标的 map
 * @param slice
 * @param fieldName
 * @return map[interface{}]interface{}
 */
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
/**
 * structToMap
 * @Description: 结构体递归转化为map
 * @param obj
 * @return map[string]interface{}
 */
func StructToMap(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldValue := reflect.ValueOf(field.Interface())
			switch field.Kind() {
			case reflect.Struct:
				result[v.Type().Field(i).Name] = StructToMap(field.Interface())
			case reflect.Slice:
				if fieldValue.Len() > 0 {
					if fieldValue.Index(0).Kind() == reflect.Struct {
						var list []map[string]interface{}
						for i := 0; i < fieldValue.Len(); i++ {
							list = append(list, StructToMap(fieldValue.Index(i).Interface()))
						}
						result[v.Type().Field(i).Name] = list
					} else {
						result[v.Type().Field(i).Name] = fieldValue.Interface()
					}
				}
			default:
				result[v.Type().Field(i).Name] = fieldValue.Interface()
			}
		}
	}
	return result
}

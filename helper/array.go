package helper

import (
	"fmt"
	"reflect"
)

// GetValueArray
// @Description: 从data中提取key的值组成数组
// @param data
// @param fieldName
// @return []interface{}
func GetValueArray(data interface{}, fieldNames ...string) []interface{} {
	var fieldValueArr []interface{}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			ele := v.Index(i).Interface()
			tmpIdArr := GetValueArray(ele, fieldNames...)
			fieldValueArr = append(fieldValueArr, tmpIdArr...)
		}
	case reflect.Struct:
		for _, fieldName := range fieldNames {
			fieldValue, err := GetStructField(data, fieldName)
			if err == nil && fieldValue != nil {
				fieldValueArr = append(fieldValueArr, fieldValue)
			}
		}
	case reflect.Map:
		for _, fieldName := range fieldNames {
			fieldValue := reflect.ValueOf(data).MapIndex(reflect.ValueOf(fieldName))
			if fieldValue.IsValid() == true {
				fieldValueArr = append(fieldValueArr, fieldValue.Interface())
			}
		}
	}
	return fieldValueArr
}

// ArrayUniqueString
// @Description:sting数组去重
// @param arr
// @return []string
func ArrayUniqueString(arr []interface{}, filterEmpty bool) []string {
	var result []string
	tempMap := map[string]byte{} // 存放不重复主键
	for _, value := range arr {
		if filterEmpty && value == "" {
			continue
		}
		l := len(tempMap)
		tempMap[value.(string)] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, value.(string))
		}
	}
	return result
}

// ArrayUnique
// @Description: 数组去重
// @param arr
// @return []string
func ArrayUnique(arr []interface{}) []any {
	var result []any
	tempMap := map[any]byte{} // 存放不重复主键
	for _, value := range arr {
		l := len(tempMap)
		tempMap[value] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, value)
		}
	}
	return result
}

// ArrayUniqueInt
// @Description: int数组去重
// @param arr
// @return []string
func ArrayUniqueInt(arr []interface{}) []int {
	var result []int
	tempMap := map[int]byte{} // 存放不重复主键
	for _, value := range arr {
		l := len(tempMap)
		tempMap[value.(int)] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, value.(int))
		}
	}
	return result
}

// ArrayUniqueInt
// @Description: int数组去重
// @param arr
// @return []string
func ArrayUniqueInt64(arr []interface{}) []int64 {
	var result []int64
	tempMap := map[int64]byte{} // 存放不重复主键
	for _, value := range arr {
		l := len(tempMap)
		tempMap[value.(int64)] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, value.(int64))
		}
	}
	return result
}

// ArrayUniqueInt
// @Description: int数组去重
// @param arr
// @return []string
func ArrayUniqueFloat64ToInt(arr []interface{}) []int {
	var result []int
	tempMap := map[float64]byte{} // 存放不重复主键
	for _, value := range arr {
		l := len(tempMap)
		tempMap[value.(float64)] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, int(value.(float64)))
		}
	}
	return result
}

// ArrayUniqueInt
// @Description: int数组去重
// @param arr
// @return []string
func ArrayUniqueFloat64(arr []interface{}) []float64 {
	var result []float64
	tempMap := map[float64]byte{} // 存放不重复主键
	for _, value := range arr {
		l := len(tempMap)
		tempMap[value.(float64)] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, value.(float64))
		}
	}
	return result
}

// GetPairs
// @Description: 返回键值对数组
// @param arr
// @param keyFieldName
// @param valueFieldName
// @return map[interface{}]interface{}
func GetPairs(arr interface{}, keyFieldName string, valueFieldName string) map[interface{}]interface{} {
	data := make(map[interface{}]interface{})
	rType := reflect.TypeOf(arr)
	if rType.Kind() != reflect.Slice {

		// log.Logger.Error("ShowPairs 函数失败，arr不是Slice类型", rType.Kind())
		// log.Logger.Error(arr)
		return data
	}
	rValue := reflect.ValueOf(arr)
	for i := 0; i < rValue.Len(); i++ {
		if reflect.TypeOf(rValue.Index(i)).Kind() == reflect.Struct {
			k, err := GetStructField(rValue.Index(i).Interface(), keyFieldName)
			if err != nil {
				fmt.Println(Now(), "提取键值对数组失败keyFieldName:", keyFieldName, err)
				// log.Logger.Error("提取键值对数组失败", keyFieldName, err)
				// log.Logger.Error(arr)
				break
			}
			v, err := GetStructField(rValue.Index(i).Interface(), valueFieldName)
			if err != nil {
				fmt.Println(Now(), "提取键值对数组失败valueFieldName:", valueFieldName, err)
				// log.Logger.Error("提取键值对数组失败", valueFieldName, err)
				// log.Logger.Error(arr)
				break
			}
			data[k] = v
		}
	}
	return data
}

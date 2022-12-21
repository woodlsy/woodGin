package helper

import (
	"reflect"
)

//
// GetValueArray
// @Description: 从data中提取key的值组成数组
// @param data
// @param fieldName
// @return []interface{}
//
func GetValueArray(data interface{}, fieldNames ...string) []interface{} {
	var fieldValueArr []interface{}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		v := reflect.ValueOf(data)
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

//
// ArrayUniqueString
// @Description:sting数组去重
// @param arr
// @return []string
//
func ArrayUniqueString(arr []interface{}) []string {
	var result []string
	tempMap := map[string]byte{} // 存放不重复主键
	for _, value := range arr {
		l := len(tempMap)
		tempMap[value.(string)] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, value.(string))
		}
	}
	return result
}

//
// ArrayUniqueInt
// @Description: int数组去重
// @param arr
// @return []string
//
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

//
// GetPairs
// @Description: 返回键值对数组
// @param arr
// @param keyFieldName
// @param valueFieldName
// @return map[interface{}]interface{}
//
func GetPairs(arr interface{}, keyFieldName string, valueFieldName string) map[interface{}]interface{} {
	data := make(map[interface{}]interface{})
	rType := reflect.TypeOf(arr)
	if rType.Kind() != reflect.Slice {

		//log.Logger.Error("ShowPairs 函数失败，arr不是Slice类型", rType.Kind())
		//log.Logger.Error(arr)
		return data
	}
	rValue := reflect.ValueOf(arr)
	for i := 0; i < rValue.Len(); i++ {
		if reflect.TypeOf(rValue.Index(i)).Kind() == reflect.Struct {
			k, err := GetStructField(rValue.Index(i).Interface(), keyFieldName)
			if err != nil {
				//log.Logger.Error("提取键值对数组失败", keyFieldName, err)
				//log.Logger.Error(arr)
				break
			}
			v, err := GetStructField(rValue.Index(i).Interface(), valueFieldName)
			if err != nil {
				//log.Logger.Error("提取键值对数组失败", valueFieldName, err)
				//log.Logger.Error(arr)
				break
			}
			data[k] = v
		}
	}
	return data
}

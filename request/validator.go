package request

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

//
// Validator
// @Description: 通过tag验证字段合法性
//             ` verify:"required;lt=2" `
//				required 必填，字符串不能为空，数字不能为0
//				lt 小于；le 小于等于；eq 等于；ne 不等于；ge 大于等于；gt 大于
//              当值是字符串时，验证的是长度；当值是数字类型时，验证的是大小；
// @param data
// @return err
//
func Validator(data interface{}) (err error) {
	//compareMap := map[string]bool{
	//	"lt": true,
	//	"le": true,
	//	"eq": true,
	//	"ne": true,
	//	"ge": true,
	//	"gt": true,
	//}
	//

	//rTyp := reflect.TypeOf(data)
	rVal := reflect.ValueOf(data) // 获取reflect.Type类型
	return verify(rVal)

}

//
// verify
// @Description: 递归验证
// @param value
// @return err
//
func verify(value reflect.Value) (err error) {
	num := value.NumField()
	for i := 0; i < num; i++ {
		field := value.Type().Field(i)
		if field.Type.Kind() == reflect.Struct {
			err = verify(value.Field(i))
			if err != nil {
				return err
			}
			continue
		}
		tag := field.Tag
		verifyRules := tag.Get("verify")
		if verifyRules == "" {
			verifyRules = tag.Get("validate")
			if verifyRules == "" {
				continue
			}
		}
		verifyRuleArr := strings.Split(verifyRules, ";")
		for r := 0; r < len(verifyRuleArr); r++ {
			if verifyRuleArr[r] == "" {
				continue
			}
			nickName := tag.Get("label")
			if nickName == "" {
				nickName = field.Name
			}
			if verifyRuleArr[r] == "required" {
				if isEmpty(value.Field(i)) {
					return errors.New(nickName + "值不能为空")
				}
			} else {
				if !compareVerify(value.Field(i), verifyRuleArr[r]) {
					return errors.New(nickName + "长度或值不在合法范围")
				}
			}
		}
	}
	return nil
}

//@function: compareVerify
//@description: 长度和数字的校验方法 根据类型自动校验
//@param: value reflect.Value, VerifyStr string
//@return: bool
func compareVerify(value reflect.Value, VerifyStr string) bool {

	switch value.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		return compare(value.Len(), VerifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), VerifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), VerifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), VerifyStr)
	default:
		return false
	}
}

//@function: isEmpty
//@description: 非空校验
//@param: value reflect.Value
//@return: bool

func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

//@function: compare
//@description: 比较函数
//@param: value interface{}, VerifyStr string
//@return: bool

func compare(value interface{}, VerifyStr string) bool {
	VerifyStrArr := strings.Split(VerifyStr, "=")
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(VerifyStrArr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Int() < VInt
		case VerifyStrArr[0] == "le":
			return val.Int() <= VInt
		case VerifyStrArr[0] == "eq":
			return val.Int() == VInt
		case VerifyStrArr[0] == "ne":
			return val.Int() != VInt
		case VerifyStrArr[0] == "ge":
			return val.Int() >= VInt
		case VerifyStrArr[0] == "gt":
			return val.Int() > VInt
		default:
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(VerifyStrArr[1])
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Uint() < uint64(VInt)
		case VerifyStrArr[0] == "le":
			return val.Uint() <= uint64(VInt)
		case VerifyStrArr[0] == "eq":
			return val.Uint() == uint64(VInt)
		case VerifyStrArr[0] == "ne":
			return val.Uint() != uint64(VInt)
		case VerifyStrArr[0] == "ge":
			return val.Uint() >= uint64(VInt)
		case VerifyStrArr[0] == "gt":
			return val.Uint() > uint64(VInt)
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStrArr[1], 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Float() < VFloat
		case VerifyStrArr[0] == "le":
			return val.Float() <= VFloat
		case VerifyStrArr[0] == "eq":
			return val.Float() == VFloat
		case VerifyStrArr[0] == "ne":
			return val.Float() != VFloat
		case VerifyStrArr[0] == "ge":
			return val.Float() >= VFloat
		case VerifyStrArr[0] == "gt":
			return val.Float() > VFloat
		default:
			return false
		}
	default:
		return false
	}
}

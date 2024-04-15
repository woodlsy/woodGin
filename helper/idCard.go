package helper

import "strconv"
import "fmt"

//
// GetGenderByIdCard
//  @Description:  根据身份证号码判断性别 0 未知 1男 2女
//  @param idCard
//  @return int
//
func GetGenderByIdCard(idCard string) int {
	gender := 0
	if len(idCard) == 0 {
		return gender
	}

	// 获取身份证号码倒数第二位
	genderDigit, err := strconv.Atoi(idCard[len(idCard)-2 : len(idCard)-1])
	if err != nil {
		return gender
	}

	// 判断性别
	if genderDigit%2 == 0 {
		gender = 2
	} else {
		gender = 1
	}

	return gender
}

//
// GetBirthdayByIdCard
//  @Description:  根据身份证号码获取生日
//  @param idCard
//  @return string
//
func GetBirthdayByIdCard(idCard string) string {
	birthday := "0001-01-01"
	if len(idCard) == 0 {
		return birthday
	}

	// 截取出生日部分
	birthdayStr := idCard[6:14]

	// 将年、月、日拼接为日期格式
	year := birthdayStr[0:4]
	month := birthdayStr[4:6]
	day := birthdayStr[6:8]

	// 组装生日字符串
	birthday = fmt.Sprintf("%s-%s-%s", year, month, day)

	return birthday
}

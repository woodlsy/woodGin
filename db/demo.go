package db

import (
	"fmt"
	"strings"
)

type tableAttribute struct {
	Collation string // 编码
	Comment   string // 备注
	//Default    interface{} // 默认值
	Extra      string // auto_increment
	Field      string
	Key        string // PRI
	Null       string
	Privileges string // 权限
	Type       string
}

func GetTableAttributes(dbName string, tableName string) []tableAttribute {

	result := make([]tableAttribute, 0)
	db[dbName].Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", tableName)).Scan(&result)
	fmt.Println("===================表结构start====================")
	for _, item := range result {
		if strings.HasPrefix(item.Type, "varchar") ||
			strings.HasPrefix(item.Type, "char") ||
			strings.HasPrefix(item.Type, "timestamp") ||
			strings.HasPrefix(item.Type, "text") ||
			strings.HasPrefix(item.Type, "date") ||
			strings.HasPrefix(item.Type, "datetime") {
			fmt.Println(UderscoreToUpperCamelCase(item.Field), "string")
		} else if strings.HasPrefix(item.Type, "int") ||
			strings.HasPrefix(item.Type, "mediumint") ||
			strings.HasPrefix(item.Type, "smallint") {
			fmt.Println(UderscoreToUpperCamelCase(item.Field), "int")
		} else if strings.HasPrefix(item.Type, "tinyint") {
			fmt.Println(UderscoreToUpperCamelCase(item.Field), "int8")
		} else {
			fmt.Println(UderscoreToUpperCamelCase(item.Field), item.Type)
		}
	}
	fmt.Println("===================表结构end====================")
	//db["at_hotel"].Raw("select id from camping limit 1").Scan(&testr)
	//db["at_hotel"].Model(&Camping{}).First(&testr)
	return result
}

func UderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}

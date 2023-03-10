package mysql

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/woodlsy/woodGin/helper"
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

func GetTableAttributes(c *gin.Context, dbName string, tableName string) []tableAttribute {

	result := make([]tableAttribute, 0)
	if _, ok := db[dbName]; !ok {
		panic(fmt.Sprintf("数据库%s未配置", dbName))
	}
	str := ""

	db[dbName].Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", tableName)).Scan(&result)
	fmt.Println("===================表结构start====================")
	str = helper.Join("", "type ", str, UderscoreToUpperCamelCase(tableName), " struct {")
	fmt.Println("type ", UderscoreToUpperCamelCase(tableName), " struct {")
	for _, item := range result {
		if strings.HasPrefix(item.Type, "varchar") ||
			strings.HasPrefix(item.Type, "char") ||
			strings.HasPrefix(item.Type, "timestamp") ||
			strings.HasPrefix(item.Type, "text") ||
			strings.HasPrefix(item.Type, "date") ||
			strings.HasPrefix(item.Type, "datetime") {
			str = helper.Join(" ", str, UderscoreToUpperCamelCase(item.Field), "string")
			fmt.Println(UderscoreToUpperCamelCase(item.Field), "string")
		} else if strings.HasPrefix(item.Type, "int") ||
			strings.HasPrefix(item.Type, "mediumint") ||
			strings.HasPrefix(item.Type, "smallint") {
			str = helper.Join(" ", str, UderscoreToUpperCamelCase(item.Field), "int")
			fmt.Println(UderscoreToUpperCamelCase(item.Field), "int")
		} else if strings.HasPrefix(item.Type, "tinyint") {
			str = helper.Join(" ", str, UderscoreToUpperCamelCase(item.Field), "int8")
			fmt.Println(UderscoreToUpperCamelCase(item.Field), "int8")
		} else {
			str = helper.Join(" ", str, UderscoreToUpperCamelCase(item.Field), item.Type)
			fmt.Println(UderscoreToUpperCamelCase(item.Field), item.Type)
		}
	}
	str = helper.Join(" ", str, "}")
	//c.HTML(200, "", str)
	fmt.Println("}")
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

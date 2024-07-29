package demo

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/woodlsy/woodGin/client/db/mysql"
	"github.com/woodlsy/woodGin/config"
	"github.com/woodlsy/woodGin/helper"
)

type Table struct {
}

func (d Table) TableAttributes(c *gin.Context) {
	params := struct {
		Base  string `form:"base"`
		Table string `form:"table"`
	}{}
	_ = c.ShouldBindQuery(&params)

	var html string
	var data gin.H
	if params.Base == "" {
		html, data = getDatabases()
	} else if params.Table == "" {
		html, data = getTables(params.Base)
	} else {
		html, data = getTableAttributes(params.Base, params.Table)
	}

	// 将HTML字符串解析成模板
	tmpl, err := template.New("index").Parse(html)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 渲染HTML模板
	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func getDatabases() (string, gin.H) {
	data := gin.H{
		"title": "数据库列表",
		"body":  config.Configs.Databases,
	}
	html := `
		<!DOCTYPE html>
        <html>
            <head>
                <title>{{ .title }}</title>
			<style>
			.main{
	text-align:center;
margin-top:100px;
}
.main a{
line-height:36px;
color:green;
}
</style>
            </head>
            <body>
<div class="main">
{{range .body}}
                <a href="?base={{.Dbname}}">{{.Dbname}}</a><br>
{{end}}
</div>
            </body>
        </html>
`
	return html, data
}

func getTables(base string) (string, gin.H) {
	var orm mysql.Orm
	orm.DatabaseName = base
	var result []map[string]interface{}
	var tables []string
	orm.Source().Raw("show tables").Scan(&result)

	for _, res := range result {
		tables = append(tables, res["Tables_in_"+base].(string))
	}
	data := gin.H{
		"title": "数据库" + base + "表列表",
		"base":  base,
		"body":  tables,
	}

	html := `
		<!DOCTYPE html>
        <html>
            <head>
                <title>{{ .title }}</title>
			<style>
			.main{
	width:200px;
margin:0 auto;
margin-top:100px;
}
.main a{
line-height:36px;
color:green;
}
</style>
            </head>
            <body>
<div class="main">
{{$tableName := .base}}
<h1>{{$tableName}}</h1>
{{range $index,$item := .body}}
                <a href="?base={{$tableName}}&table={{$item}}">{{$item}}</a><br>
{{end}}
<a href="?" style="color:blue">返回上一层</a>
</div>
            </body>
        </html>
`
	return html, data
}

type tableAttribute struct {
	Collation string // 编码
	Comment   string // 备注
	// Default    interface{} // 默认值
	Extra      string // auto_increment
	Field      string
	Key        string // PRI
	Null       string
	Privileges string // 权限
	Type       string
}

func getTableAttributes(dbName string, tableName string) (string, gin.H) {
	var orm mysql.Orm
	orm.DatabaseName = dbName
	result := make([]tableAttribute, 0)
	orm.Source().Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`", tableName)).Scan(&result)
	columns := ""

	columns = helper.Join("", "type ", columns, UderscoreToUpperCamelCase(tableName), " struct {\n")

	for _, item := range result {
		fieldName := UderscoreToUpperCamelCase(item.Field)
		fieldType := ""
		if strings.HasPrefix(item.Type, "varchar") ||
			strings.HasPrefix(item.Type, "char") ||
			strings.HasPrefix(item.Type, "timestamp") ||
			strings.HasPrefix(item.Type, "text") ||
			strings.HasPrefix(item.Type, "mediumtext") ||
			strings.HasPrefix(item.Type, "longtext") ||
			strings.HasPrefix(item.Type, "enum") ||
			strings.HasPrefix(item.Type, "date") ||
			strings.HasPrefix(item.Type, "datetime") {
			fieldType = "string"
		} else if strings.HasPrefix(item.Type, "tinyint") ||
			strings.HasPrefix(item.Type, "int") ||
			strings.HasPrefix(item.Type, "mediumint") ||
			strings.HasPrefix(item.Type, "bigint") ||
			strings.HasPrefix(item.Type, "smallint") {
			fieldType = "int64"
		} else if strings.HasPrefix(item.Type, "decimal") {
			fieldType = "float64"
		} else {
			fieldType = item.Type
		}
		columns = helper.Join(" ", columns, fieldName, fieldType, fmt.Sprintf("`json:\"%s\"`", toLowerCamelCase(fieldName)), "\n")
	}
	columns = helper.Join(" ", columns, "}")

	data := gin.H{
		"title":  "数据库" + dbName + "表" + tableName,
		"base":   tableName,
		"dbname": dbName,
		"body":   columns,
	}

	html := `
		<!DOCTYPE html>
        <html>
            <head>
                <title>{{ .title }}</title>
			<style>
			.main{
	width:600px;
margin:0 auto;
margin-top:100px;
}
.main a{
line-height:36px;
color:green;
}
</style>
            </head>
            <body>
<div class="main">
<h1>{{.base}}</h1>
                <p><textarea style="width:600px;height:600px">{{.body}}</textarea></p><br>
<a href="?base={{.dbname}}" style="color:blue">返回上一层</a>
</div>
            </body>
        </html>
`
	return html, data
}

func UderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	c := cases.Title(language.English)
	s = c.String(s)
	return strings.Replace(s, " ", "", -1)
}

func toLowerCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	// 将第一个字母转换为小写
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

package helper

import (
	"fmt"
	"strings"
)

//
// P
// @Description: 打印输出
// @param content
//
func P(content ...interface{}) {
	if len(content) == 0 {
		return
	}
	formatContent := "<!---- debug\n"
	for i := 0; i < len(content); i++ {
		formatContent = Join("", formatContent, "%+v")
	}
	formatContent = Join("", formatContent, "\ndebug ----!>\n")
	fmt.Printf(formatContent, content...)
}

//
// Join
// @Description: 拼接字符串
// @param glue 分隔符
// @param args
// @return string
//
func Join(glue string, args ...string) string {
	var build strings.Builder
	for k, s := range args {
		if k != 0 {
			build.WriteString(glue)
		}
		build.WriteString(s)
	}
	return build.String()
}

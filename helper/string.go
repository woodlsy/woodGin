package helper

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
		formatContent = Join(" ", formatContent, "%+v")
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

//
// Md5
// @Description: md5 加密
// @param str
// @return string
//
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//
// JsonEncode
// @Description: json编码数据
// @param data
// @return string
//
func JsonEncode(data interface{}) string {
	buffer := &bytes.Buffer{}
	encode := json.NewEncoder(buffer)
	encode.SetEscapeHTML(false)
	_ = encode.Encode(data)
	return string(buffer.Bytes())
}

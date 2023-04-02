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
// addslashes() 函数返回在预定义字符之前添加反斜杠的字符串。
// 预定义字符是：
// 单引号（'）
// 双引号（"）
// 反斜杠（\）
func Addslashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// stripslashes() 函数删除由 addslashes() 函数添加的反斜杠。
func Stripslashes(str string) string {
	dstRune := []rune{}
	strRune := []rune(str)
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}

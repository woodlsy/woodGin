package helper

import (
	"net/http"
	"os"
)

//
// GetFileMiMeType
// @Description: 获取文件mime类型
// @param file
// @return string
//
func GetFileMiMeType(file string) string {

	_, err := os.Stat(file)
	if err != nil {
		return ""
	}
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	// 只需要前 512 个字节就可以了
	tmpBuffer := make([]byte, 512)

	_, _ = f.Read(tmpBuffer)

	return http.DetectContentType(tmpBuffer)
}

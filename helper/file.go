package helper

import (
	"net/http"
	"os"
)

func GetFileMiMeType(file *os.File) string {
	// 只需要前 512 个字节就可以了
	tmpBuffer := make([]byte, 512)

	_, _ = file.Read(tmpBuffer)

	return http.DetectContentType(tmpBuffer)
}

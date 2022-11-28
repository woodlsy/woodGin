package config

type Log struct {
	FilePath   string `json:"filePath"` // 文件路径，不带文件名
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
	Compress   bool   `json:"compress"`
}

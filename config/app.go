package config

type App struct {
	AppName   string                 `json:"appName,omitempty"`
	PSql      bool                   `json:"PSql,omitempty"`
	Debug     bool                   `json:"debug,omitempty"`
	JwtSecret string                 `json:"jwtSecret,omitempty"`
	Mode      string                 `json:"mode,omitempty"` // 运行模式
	Custom    map[string]interface{} `json:"custom"` // 自定义配置，里面的key不能有大写
}

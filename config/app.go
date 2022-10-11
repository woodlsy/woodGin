package config

type App struct {
	AppName   string `json:"appName,omitempty"`
	PSql      bool   `json:"PSql,omitempty"`
	Debug     bool   `json:"debug,omitempty"`
	JwtSecret string `json:"jwtSecret,omitempty"`
}

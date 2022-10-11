package config

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Db       int    `json:"db"`
	Prefix   string `json:"prefix"`
	Password string `json:"password"`
}

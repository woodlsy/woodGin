package config

type Databases struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
	Prefix   string `json:"prefix"`
	Charset  string `json:"charset"`
}

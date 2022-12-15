package config

type RabbitMq struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	Vhost    string `json:"vhost"`
}

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Configs ConfigObj

func Enabled() {
	Viper("./configs/app.yml")
}

func Viper(path ...string) *viper.Viper {
	var configPath string

	configPath = path[0]
	fmt.Printf("加载配置文件%s\n", configPath)

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 实时读取配置文件
	//v.WatchConfig()
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("config file changed:", e.Name)
	//	if err = v.Unmarshal(&woodGin.Configs); err != nil {
	//		fmt.Println(err)
	//	}
	//})

	if err = v.Unmarshal(&Configs); err != nil {
		fmt.Println("配置文件解析错误", err)
	}
	return v
}

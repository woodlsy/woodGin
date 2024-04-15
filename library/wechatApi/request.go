package wechatApi

import (
	"encoding/json"
	"github.com/woodlsy/woodGin/client/curl"
	"github.com/woodlsy/woodGin/helper"
)

type Request struct {
	AppId       string
	Secret      string
	Url         string
	CustomQuery map[string]string
	CustomData  map[string]string
}

const domain = "https://api.weixin.qq.com/"

func (r Request) Get(result interface{}) {
	if r.Url == "" {
		panic("未配置请求地址")
	}

	tmpFields := []interface{}{
		helper.Join("=", "appid", r.AppId),
		helper.Join("=", "secret", r.Secret),
	}
	if len(r.CustomQuery) > 0 {
		for k, v := range r.CustomQuery {
			tmpFields = append(tmpFields, helper.Join("=", k, v))
		}
	}

	url := helper.Join("?", helper.Join("", domain, r.Url), helper.Join("&", tmpFields...))

	request := curl.Instance()
	resp := request.Get(url)
	if resp == "" {
		return
	}
	err := json.Unmarshal(request.Body, &result)
	if err != nil {
		//log.Logger.Error("json 解析天气接口数据失败", err, string(request.Body))
		return
	}
}

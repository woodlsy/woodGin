package wxa

import (
	"github.com/woodlsy/woodGin/library/wechatApi"
)

type MsgSecCheckResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		Suggest string `json:"suggest"`
		Label   int    `json:"label"`
	} `json:"result"`
}

func MsgSecCheck(accessToken string, content string, openId string) MsgSecCheckResult {
	req := wechatApi.Request{
		Url:         "wxa/msg_sec_check",
		CustomQuery: map[string]string{"access_token": accessToken},
		CustomData:  map[string]interface{}{"version": 2, "content": content, "scene": 4, "openid": openId},
	}
	var result MsgSecCheckResult
	req.Post(&result)
	return result
}

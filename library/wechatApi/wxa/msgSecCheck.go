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

func MsgSecCheck(accessToken string, version int, content string, openId string) MsgSecCheckResult {
	customData := map[string]interface{}{"content": content}
	if version == 2 {
		customData["version"] = 2
		customData["scene"] = 4
		customData["openid"] = openId
	}
	req := wechatApi.Request{
		Url:         "wxa/msg_sec_check",
		CustomQuery: map[string]string{"access_token": accessToken},
		CustomData:  customData,
	}
	var result MsgSecCheckResult
	req.Post(&result)
	return result
}

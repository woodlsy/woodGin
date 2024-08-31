package wxa

import (
	"github.com/woodlsy/woodGin/library/wechatApi"
)

type MediaCheckAsyncResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	TraceId string `json:"trace_id"`
}

func MediaCheckAsync(accessToken string, version int, url string, urlType int, openId string) MediaCheckAsyncResult {
	customData := map[string]interface{}{}
	if version == 2 {
		customData["version"] = 2
		customData["scene"] = 4
		customData["openid"] = openId
		customData["media_url"] = url
		customData["media_type"] = urlType
	} else {
		customData["media"] = url
	}
	req := wechatApi.Request{
		Url:         "wxa/img_sec_check",
		CustomQuery: map[string]string{"access_token": accessToken},
		CustomData:  customData,
	}
	var result MediaCheckAsyncResult
	req.PostLocalFile(&result, url, "media")
	return result
}

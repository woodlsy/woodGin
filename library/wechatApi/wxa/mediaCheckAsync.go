package wxa

import (
	"github.com/woodlsy/woodGin/library/wechatApi"
)

type MediaCheckAsyncResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	TraceId string `json:"trace_id"`
}

func MediaCheckAsync(accessToken string, url string, urlType int, openId string) MediaCheckAsyncResult {
	req := wechatApi.Request{
		Url:         "wxa/media_check_async",
		CustomQuery: map[string]string{"access_token": accessToken},
		CustomData:  map[string]interface{}{"version": 2, "media_url": url, "media_type": urlType, "scene": 4, "openid": openId},
	}
	var result MediaCheckAsyncResult
	req.Post(&result)
	return result
}

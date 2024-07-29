package cgiBin

import (
	"github.com/woodlsy/woodGin/library/wechatApi"
)

type TokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func Token(appId string, secret string) TokenResult {
	req := wechatApi.Request{
		AppId:       appId,
		Secret:      secret,
		Url:         "cgi-bin/token",
		CustomQuery: map[string]string{"grant_type": "client_credential"},
	}
	var result TokenResult
	req.Get(&result)
	return result
}

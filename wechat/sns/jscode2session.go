package sns

import (
	"github.com/woodlsy/woodGin/wechat"
)

type ResultJsCode2Session struct {
	Errcode    int    `json:"errcode"`
	OpenId     string `json:"openid"`
	Errmsg     string `json:"errmsg"`
	UnionId    string `json:"unionid"`
	SessionKey string `json:"session_key"`
}

func JsCode2Session(appId string, secret string, code string) ResultJsCode2Session {
	req := wechat.Request{
		AppId:       appId,
		Secret:      secret,
		Url:         "sns/jscode2session",
		CustomQuery: map[string]string{"grant_type": "authorization_code", "js_code": code},
	}
	var result ResultJsCode2Session
	req.Get(&result)
	return result
}

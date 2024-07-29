package sns

import "github.com/woodlsy/woodGin/library/wechatApi"

type JsCode2SessionResult struct {
	Errcode    int    `json:"errcode"`
	OpenId     string `json:"openid"`
	Errmsg     string `json:"errmsg"`
	UnionId    string `json:"unionid"`
	SessionKey string `json:"session_key"`
}

func JsCode2Session(appId string, secret string, code string) JsCode2SessionResult {
	req := wechatApi.Request{
		AppId:       appId,
		Secret:      secret,
		Url:         "sns/jscode2session",
		CustomQuery: map[string]string{"grant_type": "authorization_code", "js_code": code},
	}
	var result JsCode2SessionResult
	req.Get(&result)
	return result
}

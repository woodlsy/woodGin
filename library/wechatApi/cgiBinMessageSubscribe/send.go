package cgiBinMessageSubscribe

import (
	"github.com/woodlsy/woodGin/library/wechatApi"
)

type SendResult struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func Send(templateId string, toUser string, data map[string]interface{}, accessToken string, page string) SendResult {
	req := wechatApi.Request{
		Url:         "cgi-bin/message/subscribe/send",
		CustomQuery: map[string]string{"access_token": accessToken},
		CustomData: map[string]interface{}{
			"template_id": templateId,
			"touser":      toUser,
			"data":        data,
			"lang":        "zh_CN",
			"page":        page,
		},
	}
	var result SendResult
	req.Post(&result)
	return result
}

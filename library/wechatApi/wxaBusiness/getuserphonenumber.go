package wxaBusiness

import (
	"github.com/woodlsy/woodGin/library/wechatApi"
)

type GetUserPhoneNumberResult struct {
	ErrCode   int             `json:"errcode"`
	ErrMsg    string          `json:"errmsg"`
	PhoneInfo PhoneInfoResult `json:"phone_info"`
}
type PhoneInfoResult struct {
	PhoneNumber     string          `json:"phoneNumber"`     // 用户绑定的手机号（国外手机号会有区号）
	PurePhoneNumber string          `json:"purePhoneNumber"` // 没有区号的手机号
	CountryCode     string          `json:"countryCode"`     // 区号
	Watermark       WatermarkResult `json:"watermark"`
}

type WatermarkResult struct {
	Timestamp int64  `json:"timestamp"` // 用户获取手机号操作的时间戳
	AppId     string `json:"appid"`     // 小程序appid
}

func GetUserPhoneNumber(accessToken, code string) GetUserPhoneNumberResult {
	req := wechatApi.Request{
		Url:         "wxa/business/getuserphonenumber",
		CustomQuery: map[string]string{"access_token": accessToken},
		CustomData:  map[string]interface{}{"code": code},
	}
	var result GetUserPhoneNumberResult
	req.Post(&result)
	return result
}

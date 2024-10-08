package wxa

import "github.com/woodlsy/woodGin/library/wechatApi"

type LineColor struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

func GetUnlimitedQRCode(accessToken string, scene string, page string, width int, autoColor bool, lineColor *LineColor, isHyaline bool) string {
	customData := map[string]interface{}{}

	customData["scene"] = scene
	customData["auto_color"] = autoColor
	if lineColor != nil {
		customData["line_color"] = lineColor
	}
	customData["is_hyaline"] = isHyaline

	if page != "" {
		customData["page"] = page
	}
	if width != 0 {
		customData["width"] = width
	}

	req := wechatApi.Request{
		Url:         "wxa/getwxacodeunlimit",
		CustomQuery: map[string]string{"access_token": accessToken},
		CustomData:  customData,
	}
	return req.Post(nil)
}

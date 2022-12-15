package amqp

import "github.com/woodlsy/woodGin/client/amqp/rabbit"

type Amqp struct {
	Url  string
	mode string
}

var Mq Amqp

func Enabled(mode ...string) {
	if len(mode) == 0 {
		Mq.mode = "rabbit"
	} else {
		Mq.mode = mode[0]
	}
	switch Mq.mode {
	case "rabbit":
		Mq.Url = rabbit.InitUri()
	}
}

//
// PushSub
// @Description: 订阅模式推送
// @receiver a
// @param exchange
// @param queue
// @param routeKey
// @param content
// @return error
//
func (a Amqp) PushSub(exchange string, queue string, routeKey string, content map[string]interface{}) error {
	switch Mq.mode {
	case "rabbit":
		return rabbit.PushSub(a.Url, exchange, queue, routeKey, content)
	default:
		panic("无效的mq类型")
	}
	return nil
}

package rabbit

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/woodlsy/woodGin/config"
	"github.com/woodlsy/woodGin/helper"
	"github.com/woodlsy/woodGin/log"
)

func InitUri() string {
	rabbitMqConfigs := config.Configs.RabbitMq
	if rabbitMqConfigs.Host == "" || rabbitMqConfigs.UserName == "" || rabbitMqConfigs.Port == "" {
		errMsg := "请先配置mq的连接信息"
		log.Logger.Error(errMsg, rabbitMqConfigs)
		panic(errMsg)
	}
	return helper.Join("",
		"amqp://",
		rabbitMqConfigs.UserName,
		":",
		rabbitMqConfigs.Password,
		"@",
		rabbitMqConfigs.Host,
		":",
		rabbitMqConfigs.Port,
		rabbitMqConfigs.Vhost,
	)
}

// PushSub
// @Description:订阅模式生产端
// @param url
// @param exchange
// @param queue
// @param routeKey
// @param content
// @return error
func PushSub(url, exchange string, queue string, routeKey string, content string) error {
	if url == "" {
		errMsg := "未配置rabbitMq配置"
		log.Logger.Error(errMsg)
		return errors.New(errMsg)
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		errMsg := "rabbitMq连接失败"
		log.Logger.Error(errMsg, err.Error())
		return err
	}

	defer conn.Close()

	// 创建一个Channel
	channel, err := conn.Channel()
	if err != nil {
		log.Logger.Error("创建rabbit队列通道失败:", err.Error())
		return err
	}
	defer channel.Close()

	// 声明exchange
	if err := channel.ExchangeDeclare(
		exchange, //name
		"direct", //exchangeType
		true,     //durable
		false,    //auto-deleted
		false,    //internal
		false,    //noWait
		nil,      //arguments
	); err != nil {
		log.Logger.Error("Failed to declare a exchange:", err.Error())
		return err
	}
	// 声明一个queue
	if _, err := channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		log.Logger.Error("Failed to declare a queue:", err.Error())
		return err
	}
	// exchange 绑定 queue
	channel.QueueBind(queue, routeKey, exchange, false, nil)

	// 发送
	messageBody := content
	if err = channel.Publish(
		exchange, // exchange
		routeKey, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(messageBody),
			//Expiration:      "60000", // 消息过期时间
		},
	); err != nil {
		log.Logger.Error("rabbit发送消息失败:", err.Error())
		return err
	}
	return nil
}

// Consumer
// @Description: TODO 待完善 消费者
// @param url
// @param exchange
// @param queue
// @param routeKey
// @return error
func Consumer(url string, exchange string, queue string, routeKey string) error {
	// 建立连接
	conn, err := amqp.Dial(url)
	if err != nil {
		errMsg := "rabbitMq连接失败"
		log.Logger.Error(errMsg, err.Error())
		return err
	}
	defer conn.Close()
	// 启动一个通道
	channel, err := conn.Channel()
	if err != nil {
		log.Logger.Error("创建rabbit队列通道失败:", err.Error())
		return err
	}
	defer channel.Close()

	// 声明exchange
	if err := channel.ExchangeDeclare(
		exchange, //name
		"direct", //exchangeType
		true,     //durable
		false,    //auto-deleted
		false,    //internal
		false,    //noWait
		nil,      //arguments
	); err != nil {
		log.Logger.Error("Failed to declare a exchange:", err.Error())
		return err
	}

	// 声明一个队列
	q, err := channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Logger.Error("Failed to declare a queue:", err.Error())
		return err
	}
	//绑定队列到exchange中
	err = channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		routeKey,
		exchange,
		false,
		nil,
	)
	// 注册消费者
	message, err := channel.Consume(
		q.Name, // queue
		"",     // 标签
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Logger.Error("rabbit注册消费者失败:", err.Error())
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range message {
			fmt.Println(helper.Now(), d.Type)
			fmt.Println(helper.Now(), d.MessageId)
			fmt.Println(helper.Now(), d.Body)
		}
	}()
	<-forever

	return nil
}

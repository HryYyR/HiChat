package adb

import (
	"hichat_static_server/config"

	"github.com/streadway/amqp"
)

var MQc *amqp.Channel
var MQq amqp.Queue

func InitMQ() {
	conn, err := amqp.Dial(config.RabbitMQAddress)
	if err != nil {
		panic(err.Error())
	}
	channel, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}
	// 3. 声明消息要发送到的队列
	q, err := channel.QueueDeclare(
		"msgsync", // name
		false,     // durable 声明为持久队列
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		panic(err.Error())
	}

	MQq = q
	MQc = channel
}

package adb

import (
	"fmt"
	"go-websocket-server/config"

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
	// 声明消息要发送到的队列
	q, err := channel.QueueDeclare(
		"msgsync", // name
		false,     // durable
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
	fmt.Println("init rabbitmq...")
}

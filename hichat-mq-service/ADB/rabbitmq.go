package adb

import (
	"HiChat/hichat-mq-service/config"
	"github.com/streadway/amqp"
)

var MQc *amqp.Channel
var MQq amqp.Queue
var conn *amqp.Connection

func InitMQ() {
	// 尝试建立连接
	conn, err := amqp.Dial(config.RabbitMQAddress)
	if err != nil {
		panic("Failed to connect to RabbitMQ:" + err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		closeConnection()
		panic("Failed to open a channel: " + err.Error())
		return
	}

	// 声明消息要发送到的队列
	q, err := channel.QueueDeclare(
		"msgsync", // name
		false,     // durable 声明为持久队列
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		closeConnection()
		panic("Failed to declare queue: " + err.Error())
		return
	}

	MQq = q
	MQc = channel
}

// closeConnection 用于关闭AMQP连接，以避免资源泄露
func closeConnection() {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			panic("Failed to close the connection: " + err.Error())
		}
	}
}

package adb

import (
	"HiChat/hichat-mq-service/config"
	"fmt"
	"github.com/streadway/amqp"
)

var MQc *amqp.Channel
var MQq amqp.Queue
var conn *amqp.Connection

func InitMQ() {
	// 尝试建立连接
	conn, err := amqp.Dial(config.RabbitMQAddress)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s", err)
		return
	}

	channel, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %s", err)
		closeConnection()
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
		fmt.Printf("Failed to declare queue: %s", err)
		closeConnection()
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
			fmt.Printf("Failed to close the connection: %s", err)
		}
	}
}

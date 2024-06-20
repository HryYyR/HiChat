package adb

import (
	"fmt"
	"go-websocket-server/config"
	"log"

	"github.com/streadway/amqp"
)

var MqHub MqStruct

type MqStruct struct {
	MqChannel              *amqp.Channel
	NonImmediateTasksQueue amqp.Queue //用于发送非即时任务的队列(仅作为生产者,消费者位于"hichat-mq-service")
	TransmitQueue          amqp.Queue //接收消息的队列(既作为生产者,也作为消费者)
}

// PublishToNonImmediateTasksQueue 推送非即时消息到消息处理队列
func (m *MqStruct) PublishToNonImmediateTasksQueue(msg []byte) error {
	err := m.MqChannel.Publish(
		"",                            // exchange
		m.NonImmediateTasksQueue.Name, // routing key
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	if err != nil {
		return err
	}
	return nil
}

// PublishToTransmitQueue 推送所有转发消息到队列
func (m *MqStruct) PublishToTransmitQueue(msg []byte) error {

	return nil
}

// ConsumeTransmitQueue 消费转发消息
func (m *MqStruct) ConsumeTransmitQueue(msg []byte) error {

	return nil
}

func (m *MqStruct) InitMQ() {
	conn, err := amqp.Dial(config.RabbitMQAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}

	// 非即时任务消息处理队列
	q, err := channel.QueueDeclare(
		"msgsync", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	//用于转发到其他服务器的交换机
	err = channel.ExchangeDeclare(
		"msgTransmit", // 使用命名的交换器
		"fanout",      // 交换器类型
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	//转发到其他服务器的队列
	TransmitQueue, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("该服务器的MqName为:", TransmitQueue.Name)
	//队列绑定到交换机
	err = channel.QueueBind(
		m.TransmitQueue.Name, // queue name
		"",                   // routing key
		"msgTransmit",        // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.MqChannel = channel
	m.NonImmediateTasksQueue = q
	m.TransmitQueue = TransmitQueue
	fmt.Println("init rabbitmq success")
}

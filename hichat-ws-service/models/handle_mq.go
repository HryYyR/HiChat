package models

import (
	"encoding/json"
	"github.com/streadway/amqp"
	adb "go-websocket-server/ADB"
	"log"
	"strconv"
)

// RunReceiveMQMsg 消费转发队列,转发到处理中心
func RunReceiveMQMsg() {
	consume, err := adb.MqHub.MqChannel.Consume(adb.MqHub.TransmitQueue.Name, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	for delivery := range consume {
		msgbyte := delivery.Body
		//msgtype, err := strconv.Atoi(delivery.Type)
		if err != nil {
			log.Println("invaild msg type")
			continue
		}

		var basicMsg BasicMessage
		err = json.Unmarshal(msgbyte, &basicMsg)
		if err != nil {
			log.Println("Conversion user msg error: ", err)
			//continue
		}
		switch basicMsg.Type {
		case 1:
			var msg Message
			if err := json.Unmarshal(basicMsg.Data, &msg); err != nil {
				log.Printf("解析 GroupMessage 失败: %s\n", err)
				continue
			}
			ServiceCenter.Transmit <- msg

		case 2:
			var msg UserMessage
			if err := json.Unmarshal(basicMsg.Data, &msg); err != nil {
				log.Printf("解析 UserMessage 失败: %s\n", err)
				continue
			}
			ServiceCenter.Transmit <- msg
		}
	}

}

// TransmitMsg 将消息由交换机转发到其他服务器
func TransmitMsg(msgbyte []byte, types int) {
	err := adb.MqHub.MqChannel.Publish(
		"msgTransmit", // exchange
		"",            // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgbyte,
			Type:        strconv.Itoa(types),
		})
	if err != nil {
		log.Println("转发消息失败", err)
	}
}

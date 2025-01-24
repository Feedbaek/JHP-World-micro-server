package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 요청 데이터를 담을 구조체 정의
type CodeRequest struct {
	Code string `json:"code"`
}

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/" // RabbitMQ URL
	queueName   = "jhp-queue"                          // 큐 이름
)

func main() {
	// RabbitMQ 연결
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 채널 생성
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 큐 선언
	_, err = ch.QueueDeclare(
		queueName, // 큐 이름
		true,      // Durable
		false,     // Auto-deleted
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 메시지 소비
	msgs, err := ch.Consume(
		queueName, // Queue
		"",        // Consumer tag
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// 비동기로 메시지 처리
	go func() {
		for msg := range msgs {
			log.Printf("\nReceived a message: %s", msg.Body)
			res, _ := Running(string(msg.Body))
			log.Printf("\nResult: %s", res)
		}
	}()

	// 무한 대기
	select {}
}

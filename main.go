package main

import (
	"log"
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 요청 데이터를 담을 구조체 정의
type CodeRequest struct {
	Code string `json:"code"`
}

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/" // RabbitMQ URL
	subQueue = "execute-queue"                         // 읽기 큐 이름
	pubQueue = "result-queue"                          // 쓰기 큐 이름
	pubExchange = "result-exchange"                      // 익스체인지 이름
	pubRoutingKey = "result"                             // 라우팅 키
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

	// 메시지 소비
	msgs, err := ch.Consume(
		subQueue, // Queue
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
			log.Printf("\n<Received a message>\n %s", msg.Body)
			res, _ := Running(string(msg.Body))
			log.Printf("\n<Result>\n %s", res)

			// 메시지 전송
			ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
			defer cancel()

			err = ch.PublishWithContext(
				ctx,
				pubExchange,   // Exchange
				pubRoutingKey, // Routing Key
				false,     // Mandatory
				false,     // Immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(res),
				},
			)

			if err != nil {
				log.Fatalf("Failed to publish a message: %v", err)
			}
		}
	}()

	// 무한 대기
	select {}
}

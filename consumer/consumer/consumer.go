package consumer

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type handler func(message *kafka.Message)

func Consume(consumer *kafka.Consumer, topics []string, h handler) {
	run := true
	err := consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Printf("%% Unable to subscribe: %v\n", err)
	}

	for run == true {
		ev := consumer.Poll(1000)
		switch e := ev.(type) {
		case *kafka.Message:
			h(e)
		case kafka.Error:
			log.Printf("%% Error: %v\n", e)
			run = false
		}
	}
}

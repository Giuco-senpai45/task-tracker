package main

import (
	"ntf-service/consumer"
	"ntf-service/log"
	"sync"
)

func main() {
	log.Instantiate()

	topic := "tm-topic"
	kc, err := consumer.NewKafkaConsumer(topic)
	if err != nil {
		log.Error("Error creating Kafka consumer: %v", err)
	}

	var wg sync.WaitGroup
	kc.Consume(&wg)
}

package main

import (
	"ntf-service/consumer"
	"ntf-service/log"
	"sync"
)

func main() {
	log.Instantiate()

	groupId := "consumer-group"
	kc, err := consumer.NewKafkaConsumer(groupId)
	if err != nil {
		log.Error("Error creating Kafka consumer: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1) // Add a count to the WaitGroup
	consumer.Consume(kc, &wg)
	wg.Wait() // Wait for all goroutines to finish
}

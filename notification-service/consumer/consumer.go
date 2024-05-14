package consumer

import (
	"encoding/json"
	"ntf-service/log"
	"os"
	"sync"

	"github.com/IBM/sarama"
)

type MessageConsumer struct {
	Consumer sarama.Consumer
}

type Message struct {
	TaskId   int    `json:"task_id"`
	TaskName string `json:"task_name"`
	UserId   int    `json:"user_id"`
}

func NewKafkaConsumer(groupId string) (*MessageConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer([]string{"broker:29092"}, config)
	if err != nil {
		return nil, err
	}

	log.Info("Creating Kafka consumer")

	return &MessageConsumer{Consumer: consumer}, nil
}

func (c *MessageConsumer) Consume(wg *sync.WaitGroup) {
	defer wg.Done()

	partitionConsumer, err := c.Consumer.ConsumePartition("tm-topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Error("Error creating consuming partition: ", err)
	}

	signals := make(chan os.Signal, 1)

	log.Info("Starting consumer consume loop")
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var message Message
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Error("Error decoding message: ", err)
				continue
			}
			processMessage(message)

		case err := <-partitionConsumer.Errors():
			log.Error("Error: %v", err)

		case <-signals:
			break ConsumerLoop
		}
	}
}

func processMessage(message Message) {
	// Process the message. This is just a placeholder - replace with your actual processing logic.
	log.Info("Processing message: %v", message)
}

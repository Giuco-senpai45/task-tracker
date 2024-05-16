package consumer

import (
	"context"
	"encoding/json"
	"ntf-service/log"
	"os"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type MessageConsumer struct {
	Consumer sarama.Consumer
}

type Message struct {
	Type     string `json:"type"`
	TaskId   int    `json:"task_id"`
	TaskName string `json:"task_name"`
	UserId   int    `json:"user_id"`
}
type ConsumerGroupHandler struct {
	ready chan bool
}

func (consumer *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg Message
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Error("Error decoding message: %v", err)
			continue
		}
		processMessage(msg, message.Offset, message.Partition)
		session.MarkMessage(message, "")
	}
	return nil
}

func processMessage(message Message, offset int64, partition int32) {
	log.Info("Processing message: %v -  Offset: %v - Partition: %v", message, offset, partition)
}

func NewKafkaConsumer(groupId string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup([]string{"broker1:9092", "broker2:9093"}, groupId, config)
	if err != nil {
		return nil, err
	}

	log.Info("Creating Kafka consumer group")

	return consumerGroup, nil
}

func Consume(consumerGroup sarama.ConsumerGroup, wg *sync.WaitGroup) {
	defer wg.Done()

	consumer := ConsumerGroupHandler{
		ready: make(chan bool),
	}

	topic := os.Getenv("KAFKA_TOPIC")

	ctx := context.Background()
	for {
		err := consumerGroup.Consume(ctx, []string{topic}, &consumer)
		if err != nil {
			log.Error("Error from consumer: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		if ctx.Err() != nil {
			log.Error("Context error: %v", ctx.Err())
			return
		}
	}
}

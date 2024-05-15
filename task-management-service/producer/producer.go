package producer

import (
	"encoding/json"
	"tm-service/utils/log"

	"github.com/IBM/sarama"
)

type MessageProducer struct {
	Producer sarama.SyncProducer
}

type Message struct {
	TaskId   int    `json:"task_id"`
	TaskName string `json:"task_name"`
	UserId   int    `json:"user_id"`
}

func NewKafkaProducer() (*MessageProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"broker:9092"}, config)
	if err != nil {
		return nil, err
	}

	return &MessageProducer{Producer: producer}, nil
}

func GetNullMessageProducer() *MessageProducer {
	return &MessageProducer{Producer: nil}
}

func (p *MessageProducer) ProduceMessage(payload any) error {
	topic := "bt"

	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	log.Info("Produced message %v", message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return err
	}

	log.Info("Produced message to topic %s (partition %d) at offset %d", topic, partition, offset)

	return nil
}

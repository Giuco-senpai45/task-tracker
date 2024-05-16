package producer

import (
	"encoding/json"
	"fmt"
	"os"
	"tm-service/utils/log"

	"github.com/IBM/sarama"
)

type MessageProducer struct {
	Producer sarama.SyncProducer
	Topic    string
}

type Message struct {
	Type     string `json:"type"`
	TaskId   int    `json:"task_id"`
	TaskName string `json:"task_name"`
	UserId   int    `json:"user_id"`
}

func NewKafkaProducer() (*MessageProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"broker1:9092", "broker2:9093"}, config)
	if err != nil {
		return nil, err
	}

	admin, err := sarama.NewClusterAdmin([]string{"broker1:9092", "broker2:9093"}, config)
	if err != nil {
		return nil, err
	}

	topicName := os.Getenv("KAFKA_TOPIC")

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     3,
		ReplicationFactor: 1,
	}

	err = admin.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		if e, ok := err.(*sarama.TopicError); !ok || e.Err != sarama.ErrTopicAlreadyExists {
			return nil, err
		}
	}

	topics, err := admin.DescribeTopics([]string{topicName})
	if err != nil {
		return nil, err
	}

	if len(topics) > 0 {
		if len(topics[0].Partitions) < 3 {
			return nil, fmt.Errorf("topic was created with %d partitions, expected 3", len(topics[0].Partitions))
		}
	}

	return &MessageProducer{Producer: producer, Topic: topicName}, nil
}

func (p *MessageProducer) ProduceMessage(payload any) error {
	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return err
	}

	log.Info("Produced message to topic %s (partition %d) at offset %d", p.Topic, partition, offset)

	return nil
}

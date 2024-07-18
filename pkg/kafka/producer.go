package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(kafkaServerAddress string) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{kafkaServerAddress}, config)
	if err != nil {
		return Producer{}, fmt.Errorf("failed to setup producer: %w", err)
	}
	return Producer{
		producer: producer,
	}, nil
}

func (p *Producer) Emit(topic, key string, msg []byte) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(msg),
	}

	_, _, err := p.producer.SendMessage(message)
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

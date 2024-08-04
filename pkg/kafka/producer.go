package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	producer *kgo.Client
}

func NewProducer(kafkaServerAddress string) (Producer, error) {
	cl, err := kgo.NewClient(kgo.SeedBrokers(kafkaServerAddress))
	if err != nil {
		return Producer{}, fmt.Errorf("error tring to init kafka producer: %w", err)
	}

	return Producer{
		producer: cl,
	}, nil
}

func (p *Producer) Emit(topic, key string, msg []byte) error {
	record := &kgo.Record{
		Topic: topic,
		Value: msg,
		Key:   []byte(key),
	}
	if err := p.producer.ProduceSync(context.Background(), record); err.FirstErr() != nil {
		return fmt.Errorf("error tring to produce a message: %w", err.FirstErr())
	}

	return nil
}

func (p *Producer) Close() error {
	p.producer.Close()
	return nil
}

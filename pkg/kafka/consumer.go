package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Message struct {
	Key   string
	Value []byte
}

type Consumer struct {
	addr          kgo.Opt
	consumerGroup kgo.Opt
}

func NewConsumer(kafkaServerAddress, consumerGroup string) (Consumer, error) {
	c := Consumer{
		addr:          kgo.SeedBrokers(kafkaServerAddress),
		consumerGroup: kgo.ConsumerGroup(consumerGroup),
	}

	return c, nil
}

type Record struct {
	Key   string
	Value []byte
}

type Receiver func(record Record)

func (c *Consumer) Consume(ctx context.Context, topic string, receiver Receiver) error {
	consumer, err := kgo.NewClient(c.addr, c.consumerGroup, kgo.ConsumeTopics(topic))
	if err != nil {
		return fmt.Errorf("error tring to start to consume topic %s: %w", topic, err)
	}
	for {
		fetches := consumer.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			return errors.New(fmt.Sprint(errs))
		}

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			p.EachRecord(func(record *kgo.Record) {
				receiver(Record{
					Key:   string(record.Key),
					Value: record.Value,
				})
			})
		})
	}
}

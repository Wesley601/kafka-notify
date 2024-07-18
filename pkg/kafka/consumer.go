package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type Message struct {
	Key   string
	Value []byte
}

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
}

func NewConsumer(kafkaServerAddress, consumerGroup string) (Consumer, error) {
	config := sarama.NewConfig()

	cg, err := sarama.NewConsumerGroup(
		[]string{kafkaServerAddress},
		consumerGroup,
		config,
	)
	if err != nil {
		return Consumer{}, fmt.Errorf("failed to initialize consumer group: %w", err)
	}

	return Consumer{
		consumerGroup: cg,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, topic string, receive chan Message) {
	sc := saramaConsumer{receive: receive}

	for {
		err := c.consumerGroup.Consume(ctx, []string{topic}, sc)
		if err != nil {
			log.Printf("error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

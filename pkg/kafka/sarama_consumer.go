package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/wesley601/kafka-notify/pkg/models"
)

type saramaConsumer struct {
	receive chan Message
}

func (saramaConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (saramaConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (bota saramaConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		key := string(msg.Key)
		var notification models.Notification
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}

		message := Message{
			Key:   key,
			Value: msg.Value,
		}

		bota.receive <- message

		sess.MarkMessage(msg, "")
	}
	return nil
}

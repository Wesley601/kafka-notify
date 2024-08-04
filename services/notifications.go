package services

import (
	"encoding/json"
	"fmt"

	"github.com/wesley601/kafka-notify/models"
	"github.com/wesley601/kafka-notify/pkg/db"
	"github.com/wesley601/kafka-notify/pkg/kafka"
)

func ConsumeNotification(store *db.NotificationStore) kafka.Receiver {
	return func(record kafka.Record) {
		fmt.Printf("Receive the message %s\n", string(record.Value))
		var notification models.Notification

		err := json.Unmarshal(record.Value, &notification)
		if err != nil {
			fmt.Println("error tring to read the message")
			return
		}
		store.Add(record.Key, notification)
	}
}

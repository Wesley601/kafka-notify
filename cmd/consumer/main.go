package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/wesley601/kafka-notify/pkg/kafka"
	"github.com/wesley601/kafka-notify/pkg/models"
)

const (
	ConsumerGroup      = "notifications-group"
	ConsumerTopic      = "notifications"
	ConsumerPort       = ":8081"
	KafkaServerAddress = "localhost:9092"
)

var ErrNoMessagesFound = errors.New("no messages found")
var ErrNoUserFound = errors.New("no user found")

func main() {
	store := &NotificationStore{
		data: make(UserNotifications),
	}

	consumer, err := kafka.NewConsumer(KafkaServerAddress, ConsumerGroup)
	if err != nil {
		panic(err)
	}

	messages := make(chan kafka.Message)

	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Consume(ctx, ConsumerTopic, messages)
	defer cancel()

	go func() {
		for msg := range messages {
			var notification models.Notification
			err := json.Unmarshal(msg.Value, &notification)
			if err != nil {
				log.Printf("failed to unmarshal notification: %v", err)
				continue
			}
			store.Add(msg.Key, notification)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		handleNotifications(ctx, store)
	})

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 started at http://localhost%s\n", ConsumerGroup, ConsumerPort)

	if err := router.Run(ConsumerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}

func getUserIDFromRequest(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", ErrNoMessagesFound
	}
	return userID, nil
}

type UserNotifications map[string][]models.Notification

type NotificationStore struct {
	data UserNotifications
	mu   sync.RWMutex
}

func (ns *NotificationStore) Add(userID string,
	notification models.Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.data[userID] = append(ns.data[userID], notification)
}

func (ns *NotificationStore) Get(userID string) ([]models.Notification, error) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	notifications, ok := ns.data[userID]
	if !ok {
		return nil, ErrNoUserFound
	}

	return notifications, nil
}

type Consumer struct {
	store *NotificationStore
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		userID := string(msg.Key)
		var notification models.Notification
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}
		consumer.store.Add(userID, notification)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func handleNotifications(ctx *gin.Context, store *NotificationStore) {
	userID, err := getUserIDFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	notes, err := store.Get(userID)
	if err != nil {
		if errors.Is(err, ErrNoUserFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message":       fmt.Sprintf("No User found for id %s", userID),
				"notifications": []models.Notification{},
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message":       "some went wrong",
			"notifications": []models.Notification{},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": notes})
}

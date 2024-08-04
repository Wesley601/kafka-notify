package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wesley601/kafka-notify/handlers"
	"github.com/wesley601/kafka-notify/pkg/db"
	"github.com/wesley601/kafka-notify/pkg/kafka"
	"github.com/wesley601/kafka-notify/services"
)

const (
	ConsumerGroup      = "notifications-group"
	ConsumerTopic      = "notifications"
	ConsumerPort       = ":8081"
	KafkaServerAddress = "localhost:9092"
)

func main() {
	store := db.NewNotificationStore()
	consumer, err := kafka.NewConsumer(KafkaServerAddress, ConsumerGroup)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Consume(ctx, ConsumerTopic, services.ConsumeNotification(store))
	defer cancel()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/notifications/:userID", handlers.GetNotifications(store))

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 started at http://localhost%s\n", ConsumerGroup, ConsumerPort)

	if err := router.Run(ConsumerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}

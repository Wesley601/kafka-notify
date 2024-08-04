package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wesley601/kafka-notify/handlers"
	"github.com/wesley601/kafka-notify/models"
	"github.com/wesley601/kafka-notify/pkg/db"
	"github.com/wesley601/kafka-notify/pkg/kafka"
	"github.com/wesley601/kafka-notify/services"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
)

func main() {
	users := []models.User{
		{ID: 1, Name: "Emma"},
		{ID: 2, Name: "Bruno"},
		{ID: 3, Name: "Rick"},
		{ID: 4, Name: "Lena"},
	}

	producer, err := kafka.NewProducer(KafkaServerAddress)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/send", handlers.SendMessage(producer, *services.NewMessageService(producer, db.NewUserStore(users))))
	fmt.Printf("Kafka PRODUCER 📨 started at http://localhost%s\n", ProducerPort)

	if err := router.Run(ProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}

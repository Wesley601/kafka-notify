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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017/?retryWrites=true"))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	producer, err := kafka.NewProducer(KafkaServerAddress)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	userStore := db.NewUserStore(client.Database("test"))

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/send", handlers.SendMessage(producer, *services.NewMessageService(producer, userStore)))
	router.POST("/user", handlers.CreateUser(userStore))

	fmt.Printf("Kafka PRODUCER 📨 started at http://localhost%s\n", ProducerPort)

	if err := router.Run(ProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}

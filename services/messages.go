package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wesley601/kafka-notify/models"
	"github.com/wesley601/kafka-notify/pkg/kafka"
)

var KafkaTopic = "notifications"

type UserStore interface {
	ByID(int) (models.User, error)
}

type MessageService struct {
	producer  kafka.Producer
	userStore UserStore
}

func NewMessageService(producer kafka.Producer, userStore UserStore) *MessageService {
	return &MessageService{
		producer:  producer,
		userStore: userStore,
	}
}

func (ms MessageService) SendMessage(ctx *gin.Context, fromID, toID int) error {
	message := ctx.PostForm("message")

	fromUser, err := ms.userStore.ByID(fromID)
	if err != nil {
		return err
	}

	toUser, err := ms.userStore.ByID(toID)
	if err != nil {
		return err
	}

	notification := models.Notification{
		From:    fromUser,
		To:      toUser,
		Message: message,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	return ms.producer.Emit(KafkaTopic, strconv.Itoa(toUser.ID), notificationJSON)
}

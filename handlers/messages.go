package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wesley601/kafka-notify/pkg/db"
	"github.com/wesley601/kafka-notify/pkg/kafka"
	"github.com/wesley601/kafka-notify/services"
)

func SendMessage(producer kafka.Producer, messageService services.MessageService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := messageService.SendMessage(ctx.PostForm("message"), ctx.PostForm("fromID"), ctx.PostForm("toID"))
		if err != nil {
			if errors.Is(err, db.ErrUserNotFoundInProducer) {
				ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully!"})
	}
}

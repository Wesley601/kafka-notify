package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wesley601/kafka-notify/models"
	"github.com/wesley601/kafka-notify/pkg/db"
)

func GetNotifications(store *db.NotificationStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := getUserIDFromRequest(ctx)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		notes, err := store.Get(userID)
		if err != nil {
			if errors.Is(err, db.ErrNoUserFound) {
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
}

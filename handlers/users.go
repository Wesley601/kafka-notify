package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wesley601/kafka-notify/models"
	"github.com/wesley601/kafka-notify/pkg/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(store *db.UserStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := models.User{
			ID:   primitive.NewObjectID(),
			Name: ctx.PostForm("name"),
		}

		if err := store.Add(user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
	}
}

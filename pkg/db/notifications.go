package db

import (
	"context"
	"errors"

	"github.com/wesley601/kafka-notify/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoUserFound = errors.New("no user found")
)

type UserNotifications map[string][]models.Notification

type NotificationStore struct {
	collection *mongo.Collection
}

func NewNotificationStore(db *mongo.Database) *NotificationStore {
	return &NotificationStore{
		collection: db.Collection("notifications"),
	}
}

func (ns *NotificationStore) Add(notification models.Notification) error {
	notification.ID = primitive.NewObjectID()
	_, err := ns.collection.InsertOne(context.Background(), notification)
	return err
}

func (ns *NotificationStore) Get(userID string) ([]models.Notification, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	result, err := ns.collection.Find(context.Background(), bson.D{{Key: "to._id", Value: id}})
	if err != nil {
		return nil, err
	}

	var results []models.Notification
	err = result.All(context.Background(), &results)
	if err != nil {
		return nil, err
	}

	return results, err
}

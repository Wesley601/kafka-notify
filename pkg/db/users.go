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
	ErrUserNotFoundInProducer = errors.New("user not found")
)

type UserStore struct {
	collection *mongo.Collection
}

func NewUserStore(db *mongo.Database) *UserStore {
	return &UserStore{
		collection: db.Collection("users"),
	}
}

func (us *UserStore) Add(user models.User) error {
	_, err := us.collection.InsertOne(context.Background(), user)
	return err
}

func (us *UserStore) ByID(userID string) (models.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	result := us.collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, ErrUserNotFoundInProducer
		}

		return models.User{}, err
	}

	return user, nil
}

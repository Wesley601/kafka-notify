package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type Notification struct {
	ID      primitive.ObjectID `json:"-" bson:"_id"`
	From    User               `json:"from" bson:"from"`
	To      User               `json:"to" bson:"to"`
	Message string             `json:"message" bson:"message"`
}

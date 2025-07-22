package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID       bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email    string        `json:"email" bson:"email,omitempty"`
	Password bson.Binary   `json:"password" bson:"password,omitempty"`
}

package models

import "go.mongodb.org/mongo-driver/v2/bson"

type App struct {
	ID     bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name   string        `json:"name" bson:"name,omitempty"`
	Secret string        `json:"secret" bson:"secret,omitempty"`
}

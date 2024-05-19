package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,_omitempty"`
	Title       string             `json:"title" bson:"title,_omitempty"`
	Description string             `json:"description" bson:"description"`
	Author      string             `json:"author"`
	Todo_id     string             `json:"user_id"`
}

package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,_omitempty"`
	Username     string             `json:"username" bson:"username,_omitempty"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	JsonWebToken string             `json:"jsonwebtoken" bson:"jsonwebtoken"`
	User_id      string             `json:"user_id"`
}

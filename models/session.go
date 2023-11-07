package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Session struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserId    int                `bson:"user_id"`
	TokenHash string             `bson:"token_hash"`
}

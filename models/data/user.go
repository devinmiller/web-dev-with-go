package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	FirstName    string             `bson:"first_name"`
	LastName     string             `bson:"last_name"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
}

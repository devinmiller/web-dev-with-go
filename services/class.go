package services

import (
	"context"
	"fmt"

	"github.com/devinmiller/web-dev-with-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassService struct {
	client *mongo.Client
}

func NewClassService(client *mongo.Client) *ClassService {
	return &ClassService{
		client: client,
	}
}

func (s ClassService) CreateClass(ctx context.Context, className string, user *models.User) error {
	users := s.client.Database("flashy").Collection("users")

	class := models.Class{
		Id:   primitive.NewObjectID(),
		Name: className,
	}

	fmt.Println("Class to add:", class)

	filter := bson.M{"email": user.Email}
	fmt.Println("Filter:", filter)

	update := bson.M{"$push": bson.M{"classes": class}}
	fmt.Println("Update:", update)

	response, err := users.UpdateOne(ctx, filter, update)
	fmt.Println("Update response:", response)

	if err != nil {
		fmt.Println("Update error:", err)
		return err
	}

	return nil
}

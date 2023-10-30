package services

import (
	"context"

	models "github.com/devinmiller/web-dev-with-go/models/data"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	SignUp(ctx context.Context, user models.User) error
}

type userService struct {
	client *mongo.Client
}

func NewUserService(client *mongo.Client) *userService {
	return &userService{
		client: client,
	}
}

// func (s *UserService) SignIn() {

// }

func (s *userService) SignUp(ctx context.Context, user models.User) error {
	users := s.client.Database("flashy").Collection("users")

	_, err := users.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

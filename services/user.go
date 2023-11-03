package services

import (
	"context"
	"strings"

	"github.com/devinmiller/web-dev-with-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(ctx context.Context, form models.SignUpForm) error
}

type userService struct {
	client *mongo.Client
}

func NewUserService(client *mongo.Client) *userService {
	return &userService{
		client: client,
	}
}

func (s *userService) SignIn(ctx context.Context, userEmail string) (user models.User, err error) {
	users := s.client.Database("flashy").Collection("users")

	filter := bson.D{{Key: "email", Value: userEmail}}

	err = users.FindOne(ctx, filter).Decode(&user)

	return
}

func (s *userService) SignUp(ctx context.Context, form models.SignUpForm) error {
	users := s.client.Database("flashy").Collection("users")

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := models.User{
		FirstName:    form.FirstName,
		LastName:     form.LastName,
		Email:        strings.ToLower(form.Email),
		PasswordHash: string(hashedBytes),
	}

	_, err = users.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

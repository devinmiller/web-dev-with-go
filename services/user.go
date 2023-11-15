package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/devinmiller/web-dev-with-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	client *mongo.Client
}

func NewUserService(client *mongo.Client) *UserService {
	return &UserService{
		client: client,
	}
}

func (s *UserService) SignIn(ctx context.Context, form models.SignInForm) (*models.User, error) {
	users := s.client.Database("flashy").Collection("users")

	filter := bson.D{{Key: "email", Value: form.Email}}

	var user models.User
	err := users.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return nil, fmt.Errorf("authentication failure: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(form.Password))

	if err != nil {
		return nil, fmt.Errorf("authentication failure: %w", err)
	}

	return &user, nil
}

func (s *UserService) SignUp(ctx context.Context, form models.SignUpForm) error {
	users := s.client.Database("flashy").Collection("users")

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := models.User{
		Id:           primitive.NewObjectID(),
		FirstName:    form.FirstName,
		LastName:     form.LastName,
		Email:        strings.ToLower(form.Email),
		PasswordHash: string(hashedBytes),
		Classes:      make([]models.Class, 0),
	}

	_, err = users.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

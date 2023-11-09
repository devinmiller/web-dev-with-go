package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/devinmiller/web-dev-with-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionService struct {
	client *mongo.Client
}

func NewSessionService(client *mongo.Client) *SessionService {
	return &SessionService{
		client: client,
	}
}

func (s SessionService) SignIn(w http.ResponseWriter, r *http.Request, user *models.User) error {
	token, err := s.genSessionToken()
	if err != nil {
		return fmt.Errorf("session sign in: %w", err)
	}

	// Set the cookie with the generated token
	http.SetCookie(w, s.setCookie(CookieSession, token))
	// Update the user with the token hash
	user.SessionHash = s.hashSessionToken(token)

	return nil
}

func (s SessionService) SignOut(w http.ResponseWriter, r *http.Request)

const SessionTokenBytes = 32
const CookieSession = "FlashySession"

func (s SessionService) genSessionToken() (string, error) {
	b := make([]byte, SessionTokenBytes)

	_, err := rand.Read(b)

	if err != nil {
		return "", fmt.Errorf("session token: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (s SessionService) hashSessionToken(token string) string {
	// create the hash from the token
	tokenHash := sha256.Sum256([]byte(token))
	// convert to encoded string and return
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (s SessionService) setCookie(name, token string) *http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}

	return &cookie
}

func (s SessionService) updateUser(ctx context.Context, id string, hash string) error {
	users := s.client.Database("flashy").Collection("users")
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "session_hash", Value: hash},
			},
		},
	}

	result, err := users.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("set session hash: %w", err)
	}

	return nil
}

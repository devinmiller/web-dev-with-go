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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const SessionTokenBytes = 32
const CookieSession = "FlashySession"

type SessionService struct {
	client *mongo.Client
}

func NewSessionService(client *mongo.Client) *SessionService {
	return &SessionService{
		client: client,
	}
}

func (s SessionService) SignIn(w http.ResponseWriter, r *http.Request, user *models.User) error {
	token, err := s.getSessionToken()
	if err != nil {
		return fmt.Errorf("session sign in: %w", err)
	}

	// Hash the session token
	hash := s.hashSessionToken(token)
	// Set hashed session token
	err = s.setSession(r.Context(), user.Id.Hex(), hash)
	if err != nil {
		return fmt.Errorf("session sign in: %w", err)
	}

	// Set the cookie with the generated token
	s.setCookie(w, CookieSession, token, 86400)

	return nil
}

func (s SessionService) CurrentUser(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	// Get session cookie
	cookie, err := r.Cookie(CookieSession)
	if err != nil {
		return nil, fmt.Errorf("current user: %w", err)
	}
	// Hash the session token
	hash := s.hashSessionToken(cookie.Value)
	// Get user based on the session token hash
	users := s.client.Database("flashy").Collection("users")
	// Set filter to session hash
	filter := bson.D{{Key: "session_hash", Value: hash}}

	var user models.User
	err = users.FindOne(r.Context(), filter).Decode(&user)

	if err != nil {
		return nil, fmt.Errorf("current user: %w", err)
	}

	return &user, nil
}

func (s SessionService) SignOut(w http.ResponseWriter, r *http.Request) error {
	// Get the current user
	user, err := s.CurrentUser(w, r)
	if err != nil {
		return fmt.Errorf("session sign out: %w", err)
	}

	// Set empty hashed session token
	s.setSession(r.Context(), user.Id.String(), "")
	// Set cookie to expire immediately
	s.setCookie(w, CookieSession, "", -1)

	return nil
}

func (s SessionService) getSessionToken() (string, error) {
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

func (s SessionService) setCookie(w http.ResponseWriter, name string, token string, maxAge int) {
	cookie := http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
	}

	http.SetCookie(w, &cookie)
}

func (s SessionService) setSession(ctx context.Context, userId string, hash string) error {
	users := s.client.Database("flashy").Collection("users")

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("set session hash: %w", err)
	}
	// Set filter to user id
	filter := bson.D{{Key: "_id", Value: objectId}}
	// Set session has on user with id
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "session_hash", Value: hash},
			},
		},
	}

	result, err := users.UpdateOne(ctx, filter, update)

	if err != nil || result.MatchedCount == 0 {
		return fmt.Errorf("set session hash: %w", err)
	}

	return nil
}

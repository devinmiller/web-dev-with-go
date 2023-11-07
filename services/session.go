package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/devinmiller/web-dev-with-go/models"
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
	token, err := generateSessionToken()
	if err != nil {
		return fmt.Errorf("session sign in: %w", err)
	}

	cookieSession := http.Cookie{
		Name:     CookieSession,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, &cookieSession)

	return nil
}

const SessionTokenBytes = 32
const CookieSession = "FlashySession"

func generateSessionToken() (string, error) {
	b := make([]byte, SessionTokenBytes)

	_, err := rand.Read(b)

	if err != nil {
		return "", fmt.Errorf("session token: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func hashSessionToken(token string) string {
	// create the hash from the token
	tokenHash := sha256.Sum256([]byte(token))
	// convert to encoded string and return
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

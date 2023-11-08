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
	token, err := s.genSessionToken()
	if err != nil {
		return fmt.Errorf("session sign in: %w", err)
	}

	http.SetCookie(w, &cookieSession)

	user.SessionHash = s.hashSessionToken(token)

	return nil
}

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

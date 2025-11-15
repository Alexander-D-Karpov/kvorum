package security

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type Session struct {
	ID        string
	UserID    shared.ID
	ExpiresAt time.Time
	CreatedAt time.Time
}

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func NewSession(userID shared.ID, ttl time.Duration) (*Session, error) {
	token, err := GenerateSessionToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}, nil
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

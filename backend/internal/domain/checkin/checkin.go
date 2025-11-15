package checkin

import (
	"errors"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type Checkin struct {
	ID      shared.ID
	EventID shared.ID
	UserID  shared.ID
	Method  Method
	At      time.Time
}

type Method string

const (
	MethodQR     Method = "qr"
	MethodManual Method = "manual"
)

type QRToken struct {
	ID        shared.ID
	UserID    shared.ID
	EventID   shared.ID
	TokenHash []byte
	ExpiresAt time.Time
	shared.Timestamp
}

var (
	ErrInvalidQRToken   = errors.New("invalid or expired QR token")
	ErrAlreadyCheckedIn = errors.New("user already checked in")
)

func NewCheckin(eventID, userID shared.ID, method Method) *Checkin {
	return &Checkin{
		ID:      shared.NewID(),
		EventID: eventID,
		UserID:  userID,
		Method:  method,
		At:      time.Now().UTC(),
	}
}

func NewQRToken(userID, eventID shared.ID, tokenHash []byte, ttl time.Duration) *QRToken {
	return &QRToken{
		ID:        shared.NewID(),
		UserID:    userID,
		EventID:   eventID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(ttl),
		Timestamp: shared.NewTimestamp(),
	}
}

func (q *QRToken) IsExpired() bool {
	return time.Now().UTC().After(q.ExpiresAt)
}

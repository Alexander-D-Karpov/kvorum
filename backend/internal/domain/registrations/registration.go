package registrations

import (
	"encoding/json"
	"errors"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type Registration struct {
	ID      shared.ID
	EventID shared.ID
	UserID  shared.ID
	Status  Status
	Source  string
	UTM     json.RawMessage
	shared.Timestamp
}

type Status string

const (
	StatusGoing    Status = "going"
	StatusNotGoing Status = "not_going"
	StatusMaybe    Status = "maybe"
	StatusWaitlist Status = "waitlist"
)

type Waitlist struct {
	ID      shared.ID
	EventID shared.ID
	UserID  shared.ID
	shared.Timestamp
}

var (
	ErrRegistrationNotFound = errors.New("registration not found")
	ErrAlreadyRegistered    = errors.New("user already registered")
	ErrCapacityReached      = errors.New("event capacity reached")
)

func NewRegistration(eventID, userID shared.ID, source string, utm json.RawMessage) *Registration {
	return &Registration{
		ID:        shared.NewID(),
		EventID:   eventID,
		UserID:    userID,
		Status:    StatusGoing,
		Source:    source,
		UTM:       utm,
		Timestamp: shared.NewTimestamp(),
	}
}

func (r *Registration) UpdateRSVP(status Status) {
	r.Status = status
	r.Timestamp.Touch()
}

func NewWaitlistEntry(eventID, userID shared.ID) *Waitlist {
	return &Waitlist{
		ID:        shared.NewID(),
		EventID:   eventID,
		UserID:    userID,
		Timestamp: shared.NewTimestamp(),
	}
}

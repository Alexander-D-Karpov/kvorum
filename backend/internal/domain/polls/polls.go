package polls

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type Poll struct {
	ID       shared.ID
	EventID  shared.ID
	Question string
	Options  json.RawMessage
	Type     PollType
	shared.Timestamp
}

type PollType string

const (
	PollTypeSingle   PollType = "single"
	PollTypeMultiple PollType = "multiple"
	PollTypeRating   PollType = "rating"
	PollTypeNPS      PollType = "nps"
)

type Vote struct {
	ID        shared.ID
	PollID    shared.ID
	UserID    shared.ID
	OptionKey string
	CreatedAt time.Time
}

var (
	ErrPollNotFound     = errors.New("poll not found")
	ErrAlreadyVoted     = errors.New("user already voted")
	ErrInvalidOptionKey = errors.New("invalid option key")
)

func NewPoll(eventID shared.ID, question string, options json.RawMessage, pollType PollType) *Poll {
	return &Poll{
		ID:        shared.NewID(),
		EventID:   eventID,
		Question:  question,
		Options:   options,
		Type:      pollType,
		Timestamp: shared.NewTimestamp(),
	}
}

func NewVote(pollID, userID shared.ID, optionKey string) *Vote {
	return &Vote{
		ID:        shared.NewID(),
		PollID:    pollID,
		UserID:    userID,
		OptionKey: optionKey,
		CreatedAt: time.Now().UTC(),
	}
}

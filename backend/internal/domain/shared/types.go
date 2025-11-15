package shared

import (
	"time"

	"github.com/google/uuid"
)

type ID string

func NewID() ID {
	return ID(uuid.New().String())
}

func (id ID) String() string {
	return string(id)
}

type Timestamp struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTimestamp() Timestamp {
	now := time.Now().UTC()
	return Timestamp{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t *Timestamp) Touch() {
	t.UpdatedAt = time.Now().UTC()
}

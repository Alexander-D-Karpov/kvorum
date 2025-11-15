package events

import (
	"errors"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type Event struct {
	ID          shared.ID
	OwnerID     shared.ID
	Title       string
	Description string
	Visibility  Visibility
	Status      Status
	StartsAt    time.Time
	EndsAt      time.Time
	Timezone    string
	Location    string
	OnlineURL   string
	Capacity    int
	Waitlist    bool
	Settings    map[string]interface{}
	shared.Timestamp
}

type Visibility string

const (
	VisibilityPublic       Visibility = "public"
	VisibilityPrivate      Visibility = "private"
	VisibilityByLink       Visibility = "by_link"
	VisibilityByMembership Visibility = "by_membership"
	VisibilityRequest      Visibility = "request"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusCancelled Status = "cancelled"
)

type Role string

const (
	RoleOwner       Role = "owner"
	RoleOrganizer   Role = "organizer"
	RoleCoOrganizer Role = "coorganizer"
	RoleViewer      Role = "viewer"
)

type Series struct {
	ID      shared.ID
	EventID shared.ID
	RRule   string
	ExDates []time.Time
	Until   *time.Time
	shared.Timestamp
}

var (
	ErrInvalidCapacity    = errors.New("capacity must be positive or unlimited")
	ErrInvalidTimeRange   = errors.New("end time must be after start time")
	ErrCannotPublishDraft = errors.New("cannot publish event without required fields")
	ErrEventNotFound      = errors.New("event not found")
	ErrUnauthorized       = errors.New("unauthorized to perform this action")
)

func NewEvent(ownerID shared.ID, title, description string) *Event {
	return &Event{
		ID:          shared.NewID(),
		OwnerID:     ownerID,
		Title:       title,
		Description: description,
		Visibility:  VisibilityPublic,
		Status:      StatusDraft,
		Waitlist:    true,
		Settings:    make(map[string]interface{}),
		Timestamp:   shared.NewTimestamp(),
	}
}

func (e *Event) Publish() error {
	if e.Title == "" || e.StartsAt.IsZero() {
		return ErrCannotPublishDraft
	}
	if e.Status == StatusPublished {
		return nil
	}
	e.Status = StatusPublished
	e.Timestamp.Touch()
	return nil
}

func (e *Event) Cancel() {
	e.Status = StatusCancelled
	e.Timestamp.Touch()
}

func (e *Event) ValidateTimeRange() error {
	if !e.EndsAt.IsZero() && e.EndsAt.Before(e.StartsAt) {
		return ErrInvalidTimeRange
	}
	return nil
}

func (e *Event) ValidateCapacity() error {
	if e.Capacity < 0 {
		return ErrInvalidCapacity
	}
	return nil
}

func CanUserEdit(event *Event, userID shared.ID, role Role) bool {
	if event.OwnerID == userID {
		return true
	}
	return role == RoleOrganizer || role == RoleCoOrganizer
}

func CanUserPublish(event *Event, userID shared.ID, role Role) bool {
	if event.OwnerID == userID {
		return true
	}
	return role == RoleOrganizer
}

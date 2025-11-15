package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/botmax"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/checkin"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/forms"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/identity"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/registrations"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/Alexander-D-Karpov/kvorum/internal/security"
)

type Cache interface {
	SetSession(ctx context.Context, session *security.Session) error
	GetSession(ctx context.Context, sessionID string) (*security.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

type Handlers struct {
	identitySvc      *identity.Service
	eventsSvc        *events.Service
	formsSvc         *forms.Service
	registrationsSvc *registrations.Service
	checkinSvc       *checkin.Service
	pollsSvc         interface {
		CreatePoll(ctx context.Context, eventID shared.ID, question string, options json.RawMessage, pollType interface{}) (interface{}, error)
		Vote(ctx context.Context, pollID, userID shared.ID, optionKey string) error
		GetResults(ctx context.Context, pollID shared.ID) (map[string]int, error)
		GetPollsByEvent(ctx context.Context, eventID shared.ID) (interface{}, error)
	}
	calendarSvc interface {
		GenerateEventICS(ctx context.Context, eventID shared.ID) ([]byte, error)
		GenerateUserICS(ctx context.Context, userID shared.ID) ([]byte, error)
		GetGoogleCalendarLink(ctx context.Context, eventID shared.ID) (string, error)
	}
	analyticsSvc interface {
		GetEventAnalytics(ctx context.Context, eventID shared.ID, from, to time.Time) (interface{}, error)
		ExportEventAnalyticsCSV(ctx context.Context, eventID shared.ID, from, to time.Time) ([]byte, error)
	}
	botClient     *botmax.Client
	cache         Cache
	webhookSecret string
	hmacSecret    string
}

func NewHandlers(
	identitySvc *identity.Service,
	eventsSvc *events.Service,
	formsSvc *forms.Service,
	registrationsSvc *registrations.Service,
	checkinSvc *checkin.Service,
	pollsSvc interface {
		CreatePoll(ctx context.Context, eventID shared.ID, question string, options json.RawMessage, pollType interface{}) (interface{}, error)
		Vote(ctx context.Context, pollID, userID shared.ID, optionKey string) error
		GetResults(ctx context.Context, pollID shared.ID) (map[string]int, error)
		GetPollsByEvent(ctx context.Context, eventID shared.ID) (interface{}, error)
	},
	calendarSvc interface {
		GenerateEventICS(ctx context.Context, eventID shared.ID) ([]byte, error)
		GenerateUserICS(ctx context.Context, userID shared.ID) ([]byte, error)
		GetGoogleCalendarLink(ctx context.Context, eventID shared.ID) (string, error)
	},
	analyticsSvc interface {
		GetEventAnalytics(ctx context.Context, eventID shared.ID, from, to time.Time) (interface{}, error)
		ExportEventAnalyticsCSV(ctx context.Context, eventID shared.ID, from, to time.Time) ([]byte, error)
	},
	botClient *botmax.Client,
	cache Cache,
	webhookSecret string,
	hmacSecret string,
) *Handlers {
	return &Handlers{
		identitySvc:      identitySvc,
		eventsSvc:        eventsSvc,
		formsSvc:         formsSvc,
		registrationsSvc: registrationsSvc,
		checkinSvc:       checkinSvc,
		pollsSvc:         pollsSvc,
		calendarSvc:      calendarSvc,
		analyticsSvc:     analyticsSvc,
		botClient:        botClient,
		cache:            cache,
		webhookSecret:    webhookSecret,
		hmacSecret:       hmacSecret,
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

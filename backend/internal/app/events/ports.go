package events

import (
	"context"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type EventRepo interface {
	Create(ctx context.Context, event *events.Event) error
	Update(ctx context.Context, event *events.Event) error
	GetByID(ctx context.Context, id shared.ID) (*events.Event, error)
	ListPublic(ctx context.Context, limit, offset int) ([]*events.Event, error)
	Delete(ctx context.Context, id shared.ID) error
}

type SeriesRepo interface {
	Create(ctx context.Context, series *events.Series) error
	GetByEventID(ctx context.Context, eventID shared.ID) (*events.Series, error)
	Update(ctx context.Context, series *events.Series) error
	Delete(ctx context.Context, id shared.ID) error
}

type RoleRepo interface {
	Create(ctx context.Context, eventID, userID shared.ID, role events.Role) error
	GetUserRole(ctx context.Context, eventID, userID shared.ID) (events.Role, error)
	ListByEvent(ctx context.Context, eventID shared.ID) (map[shared.ID]events.Role, error)
	Delete(ctx context.Context, eventID, userID shared.ID) error
}

type Scheduler interface {
	ScheduleReminder(ctx context.Context, at time.Time, payload interface{}) (string, error)
}

type Cache interface {
	GetEventPublic(ctx context.Context, id shared.ID) (*events.Event, bool)
	SetEventPublic(ctx context.Context, id shared.ID, event *events.Event, ttl time.Duration)
	InvalidateEvent(ctx context.Context, id shared.ID)
}

type Service struct {
	eventRepo  EventRepo
	seriesRepo SeriesRepo
	roleRepo   RoleRepo
	scheduler  Scheduler
	cache      Cache
}

func NewService(eventRepo EventRepo, seriesRepo SeriesRepo, roleRepo RoleRepo, scheduler Scheduler, cache Cache) *Service {
	return &Service{
		eventRepo:  eventRepo,
		seriesRepo: seriesRepo,
		roleRepo:   roleRepo,
		scheduler:  scheduler,
		cache:      cache,
	}
}

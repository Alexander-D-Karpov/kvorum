package events

import (
	"context"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

func (s *Service) CreateEvent(ctx context.Context, userID shared.ID, title, description string) (*events.Event, error) {
	event := events.NewEvent(userID, title, description)

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	if err := s.roleRepo.Create(ctx, event.ID, userID, events.RoleOwner); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *Service) GetEvent(ctx context.Context, id shared.ID) (*events.Event, error) {
	if cached, ok := s.cache.GetEventPublic(ctx, id); ok {
		return cached, nil
	}

	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event.Status == events.StatusPublished {
		s.cache.SetEventPublic(ctx, id, event, 3*time.Minute)
	}

	return event, nil
}

func (s *Service) UpdateEvent(ctx context.Context, userID, eventID shared.ID, updates *events.Event) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	role, _ := s.roleRepo.GetUserRole(ctx, eventID, userID)
	if !events.CanUserEdit(event, userID, role) {
		return events.ErrUnauthorized
	}

	if updates.Title != "" {
		event.Title = updates.Title
	}
	if updates.Description != "" {
		event.Description = updates.Description
	}
	if !updates.StartsAt.IsZero() {
		event.StartsAt = updates.StartsAt
	}
	if !updates.EndsAt.IsZero() {
		event.EndsAt = updates.EndsAt
	}
	if updates.Location != "" {
		event.Location = updates.Location
	}
	if updates.OnlineURL != "" {
		event.OnlineURL = updates.OnlineURL
	}
	if updates.Capacity != 0 {
		event.Capacity = updates.Capacity
	}
	if updates.Visibility != "" {
		event.Visibility = updates.Visibility
	}

	if err := event.ValidateTimeRange(); err != nil {
		return err
	}
	if err := event.ValidateCapacity(); err != nil {
		return err
	}

	event.Timestamp.Touch()

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	s.cache.InvalidateEvent(ctx, eventID)
	return nil
}

func (s *Service) PublishEvent(ctx context.Context, userID, eventID shared.ID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	role, _ := s.roleRepo.GetUserRole(ctx, eventID, userID)
	if !events.CanUserPublish(event, userID, role) {
		return events.ErrUnauthorized
	}

	if err := event.Publish(); err != nil {
		return err
	}

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	if err := s.scheduleReminders(ctx, event); err != nil {
		return err
	}

	s.cache.InvalidateEvent(ctx, eventID)
	return nil
}

func (s *Service) CancelEvent(ctx context.Context, userID, eventID shared.ID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	role, _ := s.roleRepo.GetUserRole(ctx, eventID, userID)
	if !events.CanUserPublish(event, userID, role) {
		return events.ErrUnauthorized
	}

	event.Cancel()

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	s.cache.InvalidateEvent(ctx, eventID)
	return nil
}

func (s *Service) scheduleReminders(ctx context.Context, event *events.Event) error {
	reminders := []time.Duration{
		24 * time.Hour,
		3 * time.Hour,
		30 * time.Minute,
	}

	for _, d := range reminders {
		at := event.StartsAt.Add(-d)
		if at.After(time.Now()) {
			_, _ = s.scheduler.ScheduleReminder(ctx, at, map[string]interface{}{
				"event_id": event.ID,
				"type":     "reminder",
				"before":   d.String(),
			})
		}
	}

	return nil
}

func (s *Service) CreateSeries(ctx context.Context, userID, eventID shared.ID, rrule string, until *time.Time) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	role, _ := s.roleRepo.GetUserRole(ctx, eventID, userID)
	if !events.CanUserEdit(event, userID, role) {
		return events.ErrUnauthorized
	}

	series := &events.Series{
		ID:        shared.NewID(),
		EventID:   eventID,
		RRule:     rrule,
		Until:     until,
		Timestamp: shared.NewTimestamp(),
	}

	return s.seriesRepo.Create(ctx, series)
}

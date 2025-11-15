package registrations

import (
	"context"
	"encoding/json"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/registrations"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type RegistrationRepo interface {
	Create(ctx context.Context, reg *registrations.Registration) error
	GetByEventAndUser(ctx context.Context, eventID, userID shared.ID) (*registrations.Registration, error)
	Update(ctx context.Context, reg *registrations.Registration) error
	CountByEvent(ctx context.Context, eventID shared.ID, status registrations.Status) (int, error)
	ListByEvent(ctx context.Context, eventID shared.ID, statuses []registrations.Status) ([]*registrations.Registration, error)
	Delete(ctx context.Context, eventID, userID shared.ID) error
}

type WaitlistRepo interface {
	Create(ctx context.Context, entry *registrations.Waitlist) error
	GetNextByEvent(ctx context.Context, eventID shared.ID) (*registrations.Waitlist, error)
	Delete(ctx context.Context, id shared.ID) error
	CountByEvent(ctx context.Context, eventID shared.ID) (int, error)
}

type EventCapacityChecker interface {
	GetCapacity(ctx context.Context, eventID shared.ID) (int, error)
}

type Service struct {
	regRepo      RegistrationRepo
	waitlistRepo WaitlistRepo
	capacityChk  EventCapacityChecker
}

func NewService(regRepo RegistrationRepo, waitlistRepo WaitlistRepo, capacityChk EventCapacityChecker) *Service {
	return &Service{
		regRepo:      regRepo,
		waitlistRepo: waitlistRepo,
		capacityChk:  capacityChk,
	}
}

func (s *Service) Register(ctx context.Context, eventID, userID shared.ID, source string, utm json.RawMessage) (*registrations.Registration, error) {
	existing, err := s.regRepo.GetByEventAndUser(ctx, eventID, userID)
	if err == nil {
		return existing, registrations.ErrAlreadyRegistered
	}

	capacity, err := s.capacityChk.GetCapacity(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if capacity > 0 {
		count, err := s.regRepo.CountByEvent(ctx, eventID, registrations.StatusGoing)
		if err != nil {
			return nil, err
		}

		if count >= capacity {
			waitlistEntry := registrations.NewWaitlistEntry(eventID, userID)
			if err := s.waitlistRepo.Create(ctx, waitlistEntry); err != nil {
				return nil, err
			}

			reg := registrations.NewRegistration(eventID, userID, source, utm)
			reg.Status = registrations.StatusWaitlist
			if err := s.regRepo.Create(ctx, reg); err != nil {
				return nil, err
			}
			return reg, nil
		}
	}

	reg := registrations.NewRegistration(eventID, userID, source, utm)
	if err := s.regRepo.Create(ctx, reg); err != nil {
		return nil, err
	}

	return reg, nil
}

func (s *Service) UpdateRSVP(ctx context.Context, eventID, userID shared.ID, status registrations.Status) error {
	reg, err := s.regRepo.GetByEventAndUser(ctx, eventID, userID)
	if err != nil {
		return registrations.ErrRegistrationNotFound
	}

	oldStatus := reg.Status
	reg.UpdateRSVP(status)

	if err := s.regRepo.Update(ctx, reg); err != nil {
		return err
	}

	if oldStatus == registrations.StatusGoing && status != registrations.StatusGoing {
		go s.processWaitlist(context.Background(), eventID)
	}

	return nil
}

func (s *Service) CancelRegistration(ctx context.Context, eventID, userID shared.ID) error {
	if err := s.regRepo.Delete(ctx, eventID, userID); err != nil {
		return err
	}

	go s.processWaitlist(context.Background(), eventID)
	return nil
}

func (s *Service) processWaitlist(ctx context.Context, eventID shared.ID) {
	next, err := s.waitlistRepo.GetNextByEvent(ctx, eventID)
	if err != nil {
		return
	}

	reg, err := s.regRepo.GetByEventAndUser(ctx, eventID, next.UserID)
	if err != nil {
		return
	}

	reg.Status = registrations.StatusGoing
	_ = s.regRepo.Update(ctx, reg)
	_ = s.waitlistRepo.Delete(ctx, next.ID)
}

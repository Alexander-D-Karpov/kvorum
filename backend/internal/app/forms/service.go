package forms

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/forms"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type FormRepo interface {
	Create(ctx context.Context, form *forms.Form) error
	GetByID(ctx context.Context, id shared.ID) (*forms.Form, error)
	GetActiveByEvent(ctx context.Context, eventID shared.ID) (*forms.Form, error)
	Update(ctx context.Context, form *forms.Form) error
}

type ResponseRepo interface {
	Create(ctx context.Context, response *forms.Response) error
	GetByID(ctx context.Context, id shared.ID) (*forms.Response, error)
	GetByFormAndUser(ctx context.Context, formID, userID shared.ID) (*forms.Response, error)
	Update(ctx context.Context, response *forms.Response) error
	ListByForm(ctx context.Context, formID shared.ID) ([]*forms.Response, error)
}

type Cache interface {
	GetDraft(ctx context.Context, formID, userID shared.ID) (json.RawMessage, bool)
	SetDraft(ctx context.Context, formID, userID shared.ID, data json.RawMessage, ttl time.Duration)
}

type Service struct {
	formRepo     FormRepo
	responseRepo ResponseRepo
	cache        Cache
}

func NewService(formRepo FormRepo, responseRepo ResponseRepo, cache Cache) *Service {
	return &Service{
		formRepo:     formRepo,
		responseRepo: responseRepo,
		cache:        cache,
	}
}

func (s *Service) CreateForm(ctx context.Context, eventID shared.ID, schema, rules json.RawMessage) (*forms.Form, error) {
	form := forms.NewForm(eventID, schema, rules)
	if err := s.formRepo.Create(ctx, form); err != nil {
		return nil, err
	}
	return form, nil
}

func (s *Service) GetActiveForm(ctx context.Context, eventID shared.ID) (*forms.Form, error) {
	return s.formRepo.GetActiveByEvent(ctx, eventID)
}

func (s *Service) SubmitResponse(ctx context.Context, formID, userID shared.ID, answers json.RawMessage) (*forms.Response, error) {
	response, err := s.responseRepo.GetByFormAndUser(ctx, formID, userID)
	if err != nil {
		response = forms.NewResponse(formID, userID)
	}

	response.Answers = answers
	response.Submit()

	if response.Status == forms.ResponseStatusDraft {
		if err := s.responseRepo.Create(ctx, response); err != nil {
			return nil, err
		}
	} else {
		if err := s.responseRepo.Update(ctx, response); err != nil {
			return nil, err
		}
	}

	return response, nil
}

func (s *Service) SaveDraft(ctx context.Context, formID, userID shared.ID, data json.RawMessage) error {
	s.cache.SetDraft(ctx, formID, userID, data, 7*24*time.Hour)
	return nil
}

func (s *Service) GetDraft(ctx context.Context, formID, userID shared.ID) (json.RawMessage, bool) {
	return s.cache.GetDraft(ctx, formID, userID)
}

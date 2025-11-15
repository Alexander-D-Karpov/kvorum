package identity

import (
	"context"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type User struct {
	ID          shared.ID
	ProviderID  string
	Provider    string
	DisplayName string
	Email       string
	Phone       string
	Timezone    string
	Locale      string
	SavedFields map[string]interface{}
	shared.Timestamp
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id shared.ID) (*User, error)
	GetByProvider(ctx context.Context, provider, providerID string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type Service struct {
	repo UserRepo
}

func NewService(repo UserRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOrCreateUser(ctx context.Context, provider, providerID, displayName string) (*User, error) {
	user, err := s.repo.GetByProvider(ctx, provider, providerID)
	if err == nil {
		return user, nil
	}

	user = &User{
		ID:          shared.NewID(),
		Provider:    provider,
		ProviderID:  providerID,
		DisplayName: displayName,
		SavedFields: make(map[string]interface{}),
		Timestamp:   shared.NewTimestamp(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUser(ctx context.Context, id shared.ID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

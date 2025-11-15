package checkin

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/checkin"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

// CheckinRepo handles checkin records
type CheckinRepo interface {
	Create(ctx context.Context, c *checkin.Checkin) error
	GetByEventAndUser(ctx context.Context, eventID, userID shared.ID) (*checkin.Checkin, error)
	ListByEvent(ctx context.Context, eventID shared.ID) ([]*checkin.Checkin, error)
}

// QRTokenRepo handles QR tokens
type QRTokenRepo interface {
	Create(ctx context.Context, token *checkin.QRToken) error
	GetByHash(ctx context.Context, hash []byte) (*checkin.QRToken, error)
	DeleteExpired(ctx context.Context) error
}

type Service struct {
	checkinRepo CheckinRepo
	qrRepo      QRTokenRepo
	hmacSecret  []byte
}

func NewService(checkinRepo CheckinRepo, qrRepo QRTokenRepo, hmacSecret string) *Service {
	return &Service{
		checkinRepo: checkinRepo,
		qrRepo:      qrRepo,
		hmacSecret:  []byte(hmacSecret),
	}
}

func (s *Service) GenerateQRToken(ctx context.Context, userID, eventID shared.ID) (string, error) {
	data := userID.String() + ":" + eventID.String() + ":" + time.Now().Format(time.RFC3339)
	mac := hmac.New(sha256.New, s.hmacSecret)
	mac.Write([]byte(data))
	hash := mac.Sum(nil)

	token := checkin.NewQRToken(userID, eventID, hash, 24*time.Hour)
	if err := s.qrRepo.Create(ctx, token); err != nil {
		return "", err
	}

	tokenStr := base64.URLEncoding.EncodeToString(hash)
	return tokenStr, nil
}

func (s *Service) ValidateAndCheckin(ctx context.Context, tokenStr string, method checkin.Method) (*checkin.Checkin, error) {
	hash, err := base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, checkin.ErrInvalidQRToken
	}

	qrToken, err := s.qrRepo.GetByHash(ctx, hash)
	if err != nil {
		return nil, checkin.ErrInvalidQRToken
	}

	if qrToken.IsExpired() {
		return nil, checkin.ErrInvalidQRToken
	}

	c := checkin.NewCheckin(qrToken.EventID, qrToken.UserID, method)
	if err := s.checkinRepo.Create(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Service) ManualCheckin(ctx context.Context, eventID, userID shared.ID) (*checkin.Checkin, error) {
	c := checkin.NewCheckin(eventID, userID, checkin.MethodManual)
	if err := s.checkinRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

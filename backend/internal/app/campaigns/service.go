package campaigns

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type Campaign struct {
	ID         shared.ID
	EventID    shared.ID
	Name       string
	Segment    string
	Channel    string
	Content    json.RawMessage
	ScheduleAt *time.Time
	Status     string
	shared.Timestamp
}

type Delivery struct {
	ID         shared.ID
	CampaignID *shared.ID
	Channel    string
	TargetUser shared.ID
	MessageID  string
	Status     string
	Error      string
	shared.Timestamp
}

type CampaignRepo interface {
	Create(ctx context.Context, campaign *Campaign) error
	GetByID(ctx context.Context, id shared.ID) (*Campaign, error)
	ListByEvent(ctx context.Context, eventID shared.ID) ([]*Campaign, error)
	Update(ctx context.Context, campaign *Campaign) error
}

type DeliveryRepo interface {
	Create(ctx context.Context, delivery *Delivery) error
	ListByCampaign(ctx context.Context, campaignID shared.ID) ([]*Delivery, error)
}

type BotSender interface {
	SendMessage(ctx context.Context, userID shared.ID, content json.RawMessage) (string, error)
}

type Service struct {
	campaignRepo CampaignRepo
	deliveryRepo DeliveryRepo
	botSender    BotSender
}

func NewService(campaignRepo CampaignRepo, deliveryRepo DeliveryRepo, botSender BotSender) *Service {
	return &Service{
		campaignRepo: campaignRepo,
		deliveryRepo: deliveryRepo,
		botSender:    botSender,
	}
}

func (s *Service) CreateCampaign(
	ctx context.Context,
	eventID shared.ID,
	name, segment, channel, message string,
	scheduleAt *time.Time,
) (*Campaign, error) {
	content, _ := json.Marshal(map[string]string{"message": message})

	campaign := &Campaign{
		ID:         shared.NewID(),
		EventID:    eventID,
		Name:       name,
		Segment:    segment,
		Channel:    channel,
		Content:    content,
		ScheduleAt: scheduleAt,
		Status:     "pending",
		Timestamp:  shared.NewTimestamp(),
	}

	if err := s.campaignRepo.Create(ctx, campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}

func (s *Service) GetCampaigns(ctx context.Context, eventID shared.ID) ([]*Campaign, error) {
	return s.campaignRepo.ListByEvent(ctx, eventID)
}

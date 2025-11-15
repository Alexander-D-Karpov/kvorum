package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/botmax"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/hibiken/asynq"
	maxbotapi "github.com/max-messenger/max-bot-api-client-go"
)

type Event struct {
	ID          shared.ID
	Title       string
	Description string
	StartsAt    time.Time
	Timezone    string
	Location    string
	OnlineURL   string
}

type Registration struct {
	UserID shared.ID
	ChatID int64
}

type EventGetter interface {
	GetEvent(ctx context.Context, eventID shared.ID) (Event, error)
}

type RegistrationGetter interface {
	GetUserRegistrations(ctx context.Context, eventID shared.ID) ([]Registration, error)
}

type TaskHandlers struct {
	botClient   *maxbotapi.Api
	eventGetter EventGetter
	regGetter   RegistrationGetter
}

func NewTaskHandlers(botClient *maxbotapi.Api, eventGetter EventGetter, regGetter RegistrationGetter) *TaskHandlers {
	return &TaskHandlers{
		botClient:   botClient,
		eventGetter: eventGetter,
		regGetter:   regGetter,
	}
}

type ReminderPayload struct {
	EventID shared.ID     `json:"event_id"`
	Type    string        `json:"type"`
	Before  time.Duration `json:"before"`
}

func (h *TaskHandlers) HandleReminder(ctx context.Context, task *asynq.Task) error {
	var payload ReminderPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	log.Printf("Processing reminder: event=%s, before=%s", payload.EventID, payload.Before)

	if h.eventGetter == nil || h.regGetter == nil {
		log.Println("EventGetter or RegistrationGetter not set, skipping")
		return nil
	}

	event, err := h.eventGetter.GetEvent(ctx, payload.EventID)
	if err != nil {
		return fmt.Errorf("get event: %w", err)
	}

	regs, err := h.regGetter.GetUserRegistrations(ctx, payload.EventID)
	if err != nil {
		return fmt.Errorf("get registrations: %w", err)
	}

	eventForReminder := &botmax.EventForReminder{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartsAt:    event.StartsAt,
		Timezone:    event.Timezone,
		Location:    event.Location,
		OnlineURL:   event.OnlineURL,
	}

	components := botmax.BuildReminderMessageComponents(h.botClient, eventForReminder, payload.Before)

	var successCount, errorCount int
	for _, reg := range regs {
		if reg.ChatID == 0 {
			continue
		}

		msg := maxbotapi.NewMessage().
			SetChat(reg.ChatID).
			SetText(components.Text).
			SetFormat("markdown").
			AddKeyboard(components.Keyboard)

		_, err := h.botClient.Messages.Send(ctx, msg)
		if err != nil {
			log.Printf("Failed to send reminder to user %s: %v", reg.UserID, err)
			errorCount++
		} else {
			successCount++
		}
	}

	log.Printf("Reminder sent: success=%d, errors=%d", successCount, errorCount)
	return nil
}

type CampaignPayload struct {
	CampaignID string `json:"campaign_id"`
}

func (h *TaskHandlers) HandleCampaign(ctx context.Context, task *asynq.Task) error {
	var payload CampaignPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	log.Printf("Processing campaign: id=%s", payload.CampaignID)
	return nil
}

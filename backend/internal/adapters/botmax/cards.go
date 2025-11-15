package botmax

import (
	"fmt"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/registrations"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

func BuildEventCard(event *events.Event, userStatus registrations.Status) *SendMessageRequest {
	text := fmt.Sprintf("**%s**\n\n", event.Title)

	if event.Description != "" {
		text += event.Description + "\n\n"
	}

	loc, _ := time.LoadLocation(event.Timezone)
	startsAt := event.StartsAt.In(loc)
	text += fmt.Sprintf("📅 %s\n", startsAt.Format("02 Jan 2006, 15:04 MST"))

	if event.Location != "" {
		text += fmt.Sprintf("📍 %s\n", event.Location)
	}

	if event.OnlineURL != "" {
		text += fmt.Sprintf("🔗 %s\n", event.OnlineURL)
	}

	var statusEmoji string
	switch userStatus {
	case registrations.StatusGoing:
		statusEmoji = "✅ Вы идёте"
	case registrations.StatusNotGoing:
		statusEmoji = "❌ Вы не идёте"
	case registrations.StatusMaybe:
		statusEmoji = "❓ Возможно пойдёте"
	case registrations.StatusWaitlist:
		statusEmoji = "⏳ Вы в листе ожидания"
	}

	if statusEmoji != "" {
		text += fmt.Sprintf("\n%s\n", statusEmoji)
	}

	keyboard := InlineKeyboard{
		Buttons: [][]Button{
			{
				{
					Type:    "callback",
					Text:    "✅ Иду",
					Payload: FormatCallbackPayload(event.ID, "rsvp", "going"),
				},
				{
					Type:    "callback",
					Text:    "❌ Не иду",
					Payload: FormatCallbackPayload(event.ID, "rsvp", "not_going"),
				},
			},
			{
				{
					Type:    "callback",
					Text:    "❓ Возможно",
					Payload: FormatCallbackPayload(event.ID, "rsvp", "maybe"),
				},
			},
			{
				{
					Type: "link",
					Text: "ℹ️ Подробнее",
					URL:  fmt.Sprintf("https://kvorum.example.com/e/%s", event.ID),
				},
			},
		},
	}

	return &SendMessageRequest{
		Text:   text,
		Format: "markdown",
		Attachments: []Attachment{
			{
				Type:    "inline_keyboard",
				Payload: keyboard,
			},
		},
		Notify: true,
	}
}

type EventForReminder struct {
	ID          shared.ID
	Title       string
	Description string
	StartsAt    time.Time
	Timezone    string
	Location    string
	OnlineURL   string
}

func BuildReminderMessage(event *EventForReminder, before time.Duration) *SendMessageRequest {
	text := fmt.Sprintf("⏰ Напоминание: **%s**\n\n", event.Title)

	loc, _ := time.LoadLocation(event.Timezone)
	startsAt := event.StartsAt.In(loc)

	if before >= 24*time.Hour {
		text += fmt.Sprintf("Мероприятие начнётся завтра в %s\n", startsAt.Format("15:04"))
	} else if before >= time.Hour {
		hours := int(before.Hours())
		text += fmt.Sprintf("Мероприятие начнётся через %d ч. в %s\n", hours, startsAt.Format("15:04"))
	} else {
		mins := int(before.Minutes())
		text += fmt.Sprintf("Мероприятие начнётся через %d мин!\n", mins)
	}

	if event.Location != "" {
		text += fmt.Sprintf("📍 %s\n", event.Location)
	}

	keyboard := InlineKeyboard{
		Buttons: [][]Button{
			{
				{
					Type:    "callback",
					Text:    "✅ Подтвердить",
					Payload: FormatCallbackPayload(event.ID, "confirm", ""),
				},
				{
					Type:    "callback",
					Text:    "❌ Отменить",
					Payload: FormatCallbackPayload(event.ID, "cancel", ""),
				},
			},
		},
	}

	return &SendMessageRequest{
		Text:   text,
		Format: "markdown",
		Attachments: []Attachment{
			{
				Type:    "inline_keyboard",
				Payload: keyboard,
			},
		},
		Notify: true,
	}
}

func BuildWelcomeMessage(userName string) *SendMessageRequest {
	text := fmt.Sprintf("👋 Привет, %s!\n\n", userName)
	text += "Я — бот Kvorum для управления событиями.\n\n"
	text += "Я помогу тебе:\n"
	text += "• Найти интересные мероприятия\n"
	text += "• Зарегистрироваться на события\n"
	text += "• Получать напоминания\n"
	text += "• Управлять своими регистрациями\n"

	keyboard := InlineKeyboard{
		Buttons: [][]Button{
			{
				{
					Type: "link",
					Text: "🎫 Мои события",
					URL:  "https://kvorum.example.com/me",
				},
			},
		},
	}

	return &SendMessageRequest{
		Text:   text,
		Format: "markdown",
		Attachments: []Attachment{
			{
				Type:    "inline_keyboard",
				Payload: keyboard,
			},
		},
		Notify: true,
	}
}

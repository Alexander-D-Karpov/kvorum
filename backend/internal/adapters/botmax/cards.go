package botmax

import (
	"fmt"
	"time"

	domainregistrations "github.com/Alexander-D-Karpov/kvorum/internal/domain/registrations"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	maxbotapi "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type EventForCard struct {
	ID          shared.ID
	Title       string
	Description string
	StartsAt    time.Time
	Timezone    string
	Location    string
	OnlineURL   string
}

type MessageComponents struct {
	Text     string
	Keyboard *maxbotapi.Keyboard
}

func BuildEventCardComponents(api *maxbotapi.Api, event *EventForCard, userStatus domainregistrations.Status) MessageComponents {
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
	case domainregistrations.StatusGoing:
		statusEmoji = "✅ Вы идёте"
	case domainregistrations.StatusNotGoing:
		statusEmoji = "❌ Вы не идёте"
	case domainregistrations.StatusMaybe:
		statusEmoji = "❓ Возможно пойдёте"
	case domainregistrations.StatusWaitlist:
		statusEmoji = "⏳ Вы в листе ожидания"
	}

	if statusEmoji != "" {
		text += fmt.Sprintf("\n%s\n", statusEmoji)
	}

	kb := api.Messages.NewKeyboardBuilder()
	row1 := kb.AddRow()
	row1.AddCallback("✅ Иду", schemes.DEFAULT, FormatCallbackPayload(event.ID, "rsvp", "going"))
	row1.AddCallback("❌ Не иду", schemes.DEFAULT, FormatCallbackPayload(event.ID, "rsvp", "not_going"))

	row2 := kb.AddRow()
	row2.AddCallback("❓ Возможно", schemes.DEFAULT, FormatCallbackPayload(event.ID, "rsvp", "maybe"))

	row3 := kb.AddRow()
	row3.AddOpenApp("📱 Открыть мини-приложение", schemes.DEFAULT, "", fmt.Sprintf("event=%s", event.ID))

	return MessageComponents{
		Text:     text,
		Keyboard: kb,
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

func BuildReminderMessageComponents(api *maxbotapi.Api, event *EventForReminder, before time.Duration) MessageComponents {
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

	kb := api.Messages.NewKeyboardBuilder()
	row := kb.AddRow()
	row.AddCallback("✅ Подтвердить", schemes.DEFAULT, FormatCallbackPayload(event.ID, "confirm", ""))
	row.AddCallback("❌ Отменить", schemes.DEFAULT, FormatCallbackPayload(event.ID, "cancel", ""))

	row2 := kb.AddRow()
	row2.AddOpenApp("📱 Мои события", schemes.DEFAULT, "", "")

	return MessageComponents{
		Text:     text,
		Keyboard: kb,
	}
}

func BuildWelcomeMessageComponents(api *maxbotapi.Api, userName string) MessageComponents {
	text := fmt.Sprintf("👋 Привет, %s!\n\n", userName)
	text += "Я — бот Kvorum для управления событиями.\n\n"
	text += "Я помогу тебе:\n"
	text += "• Найти интересные мероприятия\n"
	text += "• Зарегистрироваться на события\n"
	text += "• Получать напоминания\n"
	text += "• Управлять своими регистрациями\n"

	kb := api.Messages.NewKeyboardBuilder()
	kb.AddRow().AddOpenApp("🎫 Мои события", schemes.DEFAULT, "", "")

	return MessageComponents{
		Text:     text,
		Keyboard: kb,
	}
}

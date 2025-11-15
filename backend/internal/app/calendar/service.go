package calendar

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type Event struct {
	ID          shared.ID
	Title       string
	Description string
	StartsAt    time.Time
	EndsAt      time.Time
	Timezone    string
	Location    string
	OnlineURL   string
}

type EventRepo interface {
	GetByID(ctx context.Context, id shared.ID) (*Event, error)
	ListByUser(ctx context.Context, userID shared.ID) ([]*Event, error)
}

type Service struct {
	eventRepo EventRepo
}

func NewService(eventRepo EventRepo) *Service {
	return &Service{eventRepo: eventRepo}
}

func (s *Service) GenerateEventICS(ctx context.Context, eventID shared.ID) ([]byte, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return generateICS([]*Event{event}), nil
}

func (s *Service) GenerateUserICS(ctx context.Context, userID shared.ID) ([]byte, error) {
	events, err := s.eventRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return generateICS(events), nil
}

func (s *Service) GetGoogleCalendarLink(ctx context.Context, eventID shared.ID) (string, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return "", err
	}

	loc, _ := time.LoadLocation(event.Timezone)
	start := event.StartsAt.In(loc).Format("20060102T150405")

	var end string
	if !event.EndsAt.IsZero() {
		end = event.EndsAt.In(loc).Format("20060102T150405")
	} else {
		end = event.StartsAt.Add(time.Hour).In(loc).Format("20060102T150405")
	}

	details := event.Description
	if event.OnlineURL != "" {
		details += "\n\nОнлайн-ссылка: " + event.OnlineURL
	}

	params := url.Values{}
	params.Set("action", "TEMPLATE")
	params.Set("text", event.Title)
	params.Set("dates", fmt.Sprintf("%s/%s", start, end))
	params.Set("details", details)
	if event.Location != "" {
		params.Set("location", event.Location)
	}

	return "https://calendar.google.com/calendar/render?" + params.Encode(), nil
}

func generateICS(events []*Event) []byte {
	var sb strings.Builder

	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString("PRODID:-//Kvorum//Event Calendar//RU\r\n")
	sb.WriteString("CALSCALE:GREGORIAN\r\n")

	for _, event := range events {
		sb.WriteString("BEGIN:VEVENT\r\n")
		sb.WriteString(fmt.Sprintf("UID:%s@kvorum.app\r\n", event.ID))
		sb.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", time.Now().UTC().Format("20060102T150405Z")))
		sb.WriteString(fmt.Sprintf("DTSTART:%s\r\n", event.StartsAt.UTC().Format("20060102T150405Z")))

		if !event.EndsAt.IsZero() {
			sb.WriteString(fmt.Sprintf("DTEND:%s\r\n", event.EndsAt.UTC().Format("20060102T150405Z")))
		}

		sb.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escapeICSText(event.Title)))

		if event.Description != "" {
			sb.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", escapeICSText(event.Description)))
		}

		if event.Location != "" {
			sb.WriteString(fmt.Sprintf("LOCATION:%s\r\n", escapeICSText(event.Location)))
		}

		if event.OnlineURL != "" {
			sb.WriteString(fmt.Sprintf("URL:%s\r\n", event.OnlineURL))
		}

		sb.WriteString("END:VEVENT\r\n")
	}

	sb.WriteString("END:VCALENDAR\r\n")

	return []byte(sb.String())
}

func escapeICSText(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, ";", "\\;")
	text = strings.ReplaceAll(text, ",", "\\,")
	text = strings.ReplaceAll(text, "\n", "\\n")
	return text
}

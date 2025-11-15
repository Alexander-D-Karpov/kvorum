package analytics

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"strconv"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

// TODO
type EventAnalytics struct {
	EventID            shared.ID        `json:"event_id"`
	PeriodFrom         time.Time        `json:"period_from"`
	PeriodTo           time.Time        `json:"period_to"`
	TotalRegistrations int              `json:"total_registrations"`
	Going              int              `json:"going"`
	NotGoing           int              `json:"not_going"`
	Maybe              int              `json:"maybe"`
	Waitlist           int              `json:"waitlist"`
	CheckedIn          int              `json:"checked_in"`
	BySource           map[string]int64 `json:"by_source"`
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetEventAnalytics(ctx context.Context, eventID shared.ID, from, to time.Time) (interface{}, error) {
	now := time.Now().UTC()
	if from.IsZero() || to.IsZero() || !from.Before(to) {
		to = now
		from = now.Add(-30 * 24 * time.Hour)
	}

	analytics := &EventAnalytics{
		EventID:            eventID,
		PeriodFrom:         from,
		PeriodTo:           to,
		TotalRegistrations: 0,
		Going:              0,
		NotGoing:           0,
		Maybe:              0,
		Waitlist:           0,
		CheckedIn:          0,
		BySource:           map[string]int64{},
	}

	return analytics, nil
}

func (s *Service) ExportEventAnalyticsCSV(ctx context.Context, eventID shared.ID, from, to time.Time) ([]byte, error) {
	data, err := s.GetEventAnalytics(ctx, eventID, from, to)
	if err != nil {
		return nil, err
	}

	analytics, ok := data.(*EventAnalytics)
	if !ok {
		return nil, errors.New("invalid analytics type")
	}

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)

	_ = writer.Write([]string{
		"event_id",
		"period_from",
		"period_to",
		"total_registrations",
		"going",
		"not_going",
		"maybe",
		"waitlist",
		"checked_in",
	})

	_ = writer.Write([]string{
		analytics.EventID.String(),
		analytics.PeriodFrom.Format(time.RFC3339),
		analytics.PeriodTo.Format(time.RFC3339),
		strconv.Itoa(analytics.TotalRegistrations),
		strconv.Itoa(analytics.Going),
		strconv.Itoa(analytics.NotGoing),
		strconv.Itoa(analytics.Maybe),
		strconv.Itoa(analytics.Waitlist),
		strconv.Itoa(analytics.CheckedIn),
	})

	writer.Flush()

	return buf.Bytes(), writer.Error()
}

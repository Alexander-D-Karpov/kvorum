package analytics

import (
	"bytes"
	"context"
	"encoding/csv"
	"strconv"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

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

type AnalyticsRepo interface {
	GetEventAnalytics(ctx context.Context, eventID shared.ID, from, to time.Time) (*EventAnalytics, error)
}

type Service struct {
	repo AnalyticsRepo
}

func NewService(repo AnalyticsRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetEventAnalytics(ctx context.Context, eventID shared.ID, from, to time.Time) (*EventAnalytics, error) {
	now := time.Now().UTC()
	if from.IsZero() || to.IsZero() || !from.Before(to) {
		to = now
		from = now.Add(-30 * 24 * time.Hour)
	}

	return s.repo.GetEventAnalytics(ctx, eventID, from, to)
}

func (s *Service) ExportEventAnalyticsCSV(ctx context.Context, eventID shared.ID, from, to time.Time) ([]byte, error) {
	analytics, err := s.GetEventAnalytics(ctx, eventID, from, to)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)

	writer.Write([]string{
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

	writer.Write([]string{
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

	writer.Write([]string{})
	writer.Write([]string{"Source", "Count"})

	for source, count := range analytics.BySource {
		writer.Write([]string{source, strconv.FormatInt(count, 10)})
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

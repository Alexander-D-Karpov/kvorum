package repo

import (
	"context"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/app/analytics"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type AnalyticsRepo struct {
	db *DB
}

func NewAnalyticsRepo(db *DB) *AnalyticsRepo {
	return &AnalyticsRepo{db: db}
}

func (r *AnalyticsRepo) GetEventAnalytics(ctx context.Context, eventID shared.ID, from, to time.Time) (*analytics.EventAnalytics, error) {
	query := `
        WITH stats AS (
            SELECT 
                COUNT(*) FILTER (WHERE status = 'going') as going_count,
                COUNT(*) FILTER (WHERE status = 'not_going') as not_going_count,
                COUNT(*) FILTER (WHERE status = 'maybe') as maybe_count,
                COUNT(*) FILTER (WHERE status = 'waitlist') as waitlist_count,
                COUNT(*) as total_count
            FROM registrations
            WHERE event_id = $1
              AND created_at >= $2
              AND created_at <= $3
        ),
        checkins AS (
            SELECT COUNT(*) as checkin_count
            FROM checkins
            WHERE event_id = $1
              AND at >= $2
              AND at <= $3
        ),
        sources AS (
            SELECT 
                COALESCE(source, 'unknown') as source_name,
                COUNT(*) as source_count
            FROM registrations
            WHERE event_id = $1
              AND created_at >= $2
              AND created_at <= $3
            GROUP BY source
        )
        SELECT 
            s.total_count,
            s.going_count,
            s.not_going_count,
            s.maybe_count,
            s.waitlist_count,
            COALESCE(c.checkin_count, 0) as checkin_count
        FROM stats s
        CROSS JOIN checkins c
    `

	var result analytics.EventAnalytics
	result.EventID = eventID
	result.PeriodFrom = from
	result.PeriodTo = to
	result.BySource = make(map[string]int64)

	err := r.db.pool.QueryRow(ctx, query, eventID, from, to).Scan(
		&result.TotalRegistrations,
		&result.Going,
		&result.NotGoing,
		&result.Maybe,
		&result.Waitlist,
		&result.CheckedIn,
	)
	if err != nil {
		return nil, err
	}

	sourcesQuery := `
        SELECT 
            COALESCE(source, 'unknown') as source_name,
            COUNT(*) as source_count
        FROM registrations
        WHERE event_id = $1
          AND created_at >= $2
          AND created_at <= $3
        GROUP BY source
    `

	rows, err := r.db.pool.Query(ctx, sourcesQuery, eventID, from, to)
	if err != nil {
		return &result, nil
	}
	defer rows.Close()

	for rows.Next() {
		var source string
		var count int64
		if err := rows.Scan(&source, &count); err != nil {
			continue
		}
		result.BySource[source] = count
	}

	return &result, nil
}

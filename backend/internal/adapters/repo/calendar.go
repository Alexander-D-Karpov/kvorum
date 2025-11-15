package repo

import (
	"context"

	"github.com/Alexander-D-Karpov/kvorum/internal/app/calendar"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type CalendarEventRepo struct {
	db *DB
}

func NewCalendarEventRepo(db *DB) *CalendarEventRepo {
	return &CalendarEventRepo{db: db}
}

func (r *CalendarEventRepo) GetByID(ctx context.Context, id shared.ID) (*calendar.Event, error) {
	query := `
		SELECT id, title, description, starts_at, ends_at, tz, location, online_url
		FROM events
		WHERE id = $1
	`

	var event calendar.Event
	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&event.ID, &event.Title, &event.Description,
		&event.StartsAt, &event.EndsAt, &event.Timezone,
		&event.Location, &event.OnlineURL,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *CalendarEventRepo) ListByUser(ctx context.Context, userID shared.ID) ([]*calendar.Event, error) {
	query := `
		SELECT e.id, e.title, e.description, e.starts_at, e.ends_at, e.tz, e.location, e.online_url
		FROM events e
		JOIN registrations r ON r.event_id = e.id
		WHERE r.user_id = $1 AND r.status = 'going'
		ORDER BY e.starts_at DESC
	`

	rows, err := r.db.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*calendar.Event
	for rows.Next() {
		var event calendar.Event
		err := rows.Scan(
			&event.ID, &event.Title, &event.Description,
			&event.StartsAt, &event.EndsAt, &event.Timezone,
			&event.Location, &event.OnlineURL,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &event)
	}

	return result, rows.Err()
}

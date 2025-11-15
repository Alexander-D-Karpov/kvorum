package repo

import (
	"context"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type SeriesRepo struct {
	db *DB
}

func NewSeriesRepo(db *DB) *SeriesRepo {
	return &SeriesRepo{db: db}
}

func (r *SeriesRepo) Create(ctx context.Context, series *events.Series) error {
	query := `
		INSERT INTO event_series (id, event_id, rrule, exdates, until, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.pool.Exec(ctx, query,
		series.ID, series.EventID, series.RRule, series.ExDates,
		series.Until, series.CreatedAt, series.UpdatedAt,
	)
	return err
}

func (r *SeriesRepo) GetByEventID(ctx context.Context, eventID shared.ID) (*events.Series, error) {
	query := `
		SELECT id, event_id, rrule, exdates, until, created_at, updated_at
		FROM event_series
		WHERE event_id = $1
	`

	var series events.Series
	err := r.db.pool.QueryRow(ctx, query, eventID).Scan(
		&series.ID, &series.EventID, &series.RRule, &series.ExDates,
		&series.Until, &series.CreatedAt, &series.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &series, nil
}

func (r *SeriesRepo) Update(ctx context.Context, series *events.Series) error {
	query := `
		UPDATE event_series
		SET rrule = $2, exdates = $3, until = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.pool.Exec(ctx, query,
		series.ID, series.RRule, series.ExDates, series.Until, series.UpdatedAt,
	)
	return err
}

func (r *SeriesRepo) Delete(ctx context.Context, id shared.ID) error {
	query := `DELETE FROM event_series WHERE id = $1`
	_, err := r.db.pool.Exec(ctx, query, id)
	return err
}

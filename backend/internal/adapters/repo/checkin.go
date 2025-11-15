package repo

import (
	"context"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/checkin"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type CheckinRepo struct {
	db *DB
}

func NewCheckinRepo(db *DB) *CheckinRepo {
	return &CheckinRepo{db: db}
}

func (r *CheckinRepo) Create(ctx context.Context, c *checkin.Checkin) error {
	query := `
		INSERT INTO checkins (id, event_id, user_id, method, at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.pool.Exec(ctx, query, c.ID, c.EventID, c.UserID, c.Method, c.At)
	return err
}

func (r *CheckinRepo) GetByEventAndUser(ctx context.Context, eventID, userID shared.ID) (*checkin.Checkin, error) {
	query := `
		SELECT id, event_id, user_id, method, at
		FROM checkins
		WHERE event_id = $1 AND user_id = $2
		ORDER BY at DESC
		LIMIT 1
	`

	var c checkin.Checkin
	err := r.db.pool.QueryRow(ctx, query, eventID, userID).Scan(
		&c.ID, &c.EventID, &c.UserID, &c.Method, &c.At,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CheckinRepo) ListByEvent(ctx context.Context, eventID shared.ID) ([]*checkin.Checkin, error) {
	query := `
		SELECT id, event_id, user_id, method, at
		FROM checkins
		WHERE event_id = $1
		ORDER BY at DESC
	`

	rows, err := r.db.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*checkin.Checkin
	for rows.Next() {
		var c checkin.Checkin
		err := rows.Scan(&c.ID, &c.EventID, &c.UserID, &c.Method, &c.At)
		if err != nil {
			return nil, err
		}
		result = append(result, &c)
	}

	return result, rows.Err()
}

type QRTokenRepo struct {
	db *DB
}

func NewQRTokenRepo(db *DB) *QRTokenRepo {
	return &QRTokenRepo{db: db}
}

func (r *QRTokenRepo) Create(ctx context.Context, token *checkin.QRToken) error {
	query := `
		INSERT INTO qr_tokens (id, user_id, event_id, token_hash, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.pool.Exec(ctx, query,
		token.ID, token.UserID, token.EventID, token.TokenHash,
		token.ExpiresAt, token.CreatedAt, token.UpdatedAt,
	)
	return err
}

func (r *QRTokenRepo) GetByHash(ctx context.Context, hash []byte) (*checkin.QRToken, error) {
	query := `
		SELECT id, user_id, event_id, token_hash, expires_at, created_at, updated_at
		FROM qr_tokens
		WHERE token_hash = $1
	`

	var token checkin.QRToken
	err := r.db.pool.QueryRow(ctx, query, hash).Scan(
		&token.ID, &token.UserID, &token.EventID, &token.TokenHash,
		&token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, checkin.ErrInvalidQRToken
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *QRTokenRepo) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM qr_tokens WHERE expires_at < NOW()`
	_, err := r.db.pool.Exec(ctx, query)
	return err
}

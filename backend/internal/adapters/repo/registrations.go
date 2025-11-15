package repo

import (
	"context"
	"fmt"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/registrations"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type RegistrationRepo struct {
	db *DB
}

func NewRegistrationRepo(db *DB) *RegistrationRepo {
	return &RegistrationRepo{db: db}
}

func (r *RegistrationRepo) Create(ctx context.Context, reg *registrations.Registration) error {
	query := `
		INSERT INTO registrations (id, event_id, user_id, status, source, utm, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.pool.Exec(ctx, query,
		reg.ID, reg.EventID, reg.UserID, reg.Status, reg.Source, reg.UTM, reg.CreatedAt, reg.UpdatedAt,
	)
	return err
}

func (r *RegistrationRepo) GetByEventAndUser(ctx context.Context, eventID, userID shared.ID) (*registrations.Registration, error) {
	query := `
		SELECT id, event_id, user_id, status, source, utm, created_at, updated_at
		FROM registrations
		WHERE event_id = $1 AND user_id = $2
	`

	var reg registrations.Registration
	err := r.db.pool.QueryRow(ctx, query, eventID, userID).Scan(
		&reg.ID, &reg.EventID, &reg.UserID, &reg.Status, &reg.Source, &reg.UTM, &reg.CreatedAt, &reg.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, registrations.ErrRegistrationNotFound
	}
	if err != nil {
		return nil, err
	}

	return &reg, nil
}

func (r *RegistrationRepo) Update(ctx context.Context, reg *registrations.Registration) error {
	query := `
		UPDATE registrations
		SET status = $3, updated_at = $4
		WHERE event_id = $1 AND user_id = $2
	`
	_, err := r.db.pool.Exec(ctx, query, reg.EventID, reg.UserID, reg.Status, reg.UpdatedAt)
	return err
}

func (r *RegistrationRepo) CountByEvent(ctx context.Context, eventID shared.ID, status registrations.Status) (int, error) {
	query := `SELECT COUNT(*) FROM registrations WHERE event_id = $1 AND status = $2`
	var count int
	err := r.db.pool.QueryRow(ctx, query, eventID, status).Scan(&count)
	return count, err
}

func (r *RegistrationRepo) ListByEvent(ctx context.Context, eventID shared.ID, statuses []registrations.Status) ([]*registrations.Registration, error) {
	query := `
		SELECT id, event_id, user_id, status, source, utm, created_at, updated_at
		FROM registrations
		WHERE event_id = $1 AND status = ANY($2)
		ORDER BY created_at DESC
	`

	rows, err := r.db.pool.Query(ctx, query, eventID, statuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*registrations.Registration
	for rows.Next() {
		var reg registrations.Registration
		err := rows.Scan(
			&reg.ID, &reg.EventID, &reg.UserID, &reg.Status, &reg.Source, &reg.UTM, &reg.CreatedAt, &reg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &reg)
	}

	return result, rows.Err()
}

func (r *RegistrationRepo) Delete(ctx context.Context, eventID, userID shared.ID) error {
	query := `DELETE FROM registrations WHERE event_id = $1 AND user_id = $2`
	_, err := r.db.pool.Exec(ctx, query, eventID, userID)
	return err
}

type WaitlistRepo struct {
	db *DB
}

func NewWaitlistRepo(db *DB) *WaitlistRepo {
	return &WaitlistRepo{db: db}
}

func (r *WaitlistRepo) Create(ctx context.Context, entry *registrations.Waitlist) error {
	query := `
		INSERT INTO waitlist (id, event_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.pool.Exec(ctx, query, entry.ID, entry.EventID, entry.UserID, entry.CreatedAt)
	return err
}

func (r *WaitlistRepo) GetNextByEvent(ctx context.Context, eventID shared.ID) (*registrations.Waitlist, error) {
	query := `
		SELECT id, event_id, user_id, created_at
		FROM waitlist
		WHERE event_id = $1
		ORDER BY created_at ASC
		LIMIT 1
	`

	var entry registrations.Waitlist
	err := r.db.pool.QueryRow(ctx, query, eventID).Scan(
		&entry.ID, &entry.EventID, &entry.UserID, &entry.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (r *WaitlistRepo) Delete(ctx context.Context, id shared.ID) error {
	query := `DELETE FROM waitlist WHERE id = $1`
	_, err := r.db.pool.Exec(ctx, query, id)
	return err
}

func (r *WaitlistRepo) CountByEvent(ctx context.Context, eventID shared.ID) (int, error) {
	query := `SELECT COUNT(*) FROM waitlist WHERE event_id = $1`
	var count int
	err := r.db.pool.QueryRow(ctx, query, eventID).Scan(&count)
	return count, err
}

func (r *RegistrationRepo) GetByEventWithChatIDs(ctx context.Context, eventID shared.ID) ([]struct {
	UserID shared.ID
	ChatID int64
}, error) {
	query := `
		SELECT r.user_id, ui.provider_user_id
		FROM registrations r
		JOIN user_identities ui ON ui.user_id = r.user_id AND ui.provider = 'max'
		WHERE r.event_id = $1 AND r.status = 'going'
	`

	rows, err := r.db.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []struct {
		UserID shared.ID
		ChatID int64
	}

	for rows.Next() {
		var userID shared.ID
		var providerUserID string

		if err := rows.Scan(&userID, &providerUserID); err != nil {
			return nil, err
		}

		var chatID int64
		fmt.Sscanf(providerUserID, "%d", &chatID)

		result = append(result, struct {
			UserID shared.ID
			ChatID int64
		}{
			UserID: userID,
			ChatID: chatID,
		})
	}

	return result, rows.Err()
}

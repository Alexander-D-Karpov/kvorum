package repo

import (
	"context"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type EventRepo struct {
	db *DB
}

func NewEventRepo(db *DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Create(ctx context.Context, event *events.Event) error {
	query := `
		INSERT INTO events (
			id, owner_id, title, description, visibility, status,
			starts_at, ends_at, tz, location, online_url,
			capacity, waitlist_enabled, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`
	_, err := r.db.pool.Exec(ctx, query,
		event.ID, event.OwnerID, event.Title, event.Description,
		event.Visibility, event.Status, event.StartsAt, event.EndsAt,
		event.Timezone, event.Location, event.OnlineURL, event.Capacity,
		event.Waitlist, event.CreatedAt, event.UpdatedAt,
	)
	return err
}

func (r *EventRepo) Update(ctx context.Context, event *events.Event) error {
	query := `
		UPDATE events SET
			title = $2, description = $3, visibility = $4, status = $5,
			starts_at = $6, ends_at = $7, tz = $8, location = $9,
			online_url = $10, capacity = $11, waitlist_enabled = $12,
			updated_at = $13
		WHERE id = $1
	`
	_, err := r.db.pool.Exec(ctx, query,
		event.ID, event.Title, event.Description, event.Visibility,
		event.Status, event.StartsAt, event.EndsAt, event.Timezone,
		event.Location, event.OnlineURL, event.Capacity, event.Waitlist,
		event.UpdatedAt,
	)
	return err
}

func (r *EventRepo) GetByID(ctx context.Context, id shared.ID) (*events.Event, error) {
	query := `
		SELECT id, owner_id, title, description, visibility, status,
		       starts_at, ends_at, tz, location, online_url,
		       capacity, waitlist_enabled, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	var event events.Event
	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&event.ID, &event.OwnerID, &event.Title, &event.Description,
		&event.Visibility, &event.Status, &event.StartsAt, &event.EndsAt,
		&event.Timezone, &event.Location, &event.OnlineURL,
		&event.Capacity, &event.Waitlist, &event.CreatedAt, &event.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, events.ErrEventNotFound
	}
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventRepo) GetCapacity(ctx context.Context, eventID shared.ID) (int, error) {
	query := `SELECT capacity FROM events WHERE id = $1`

	var capacity int
	err := r.db.pool.QueryRow(ctx, query, eventID).Scan(&capacity)
	if err == pgx.ErrNoRows {
		return 0, events.ErrEventNotFound
	}
	if err != nil {
		return 0, err
	}

	return capacity, nil
}

func (r *EventRepo) ListPublic(ctx context.Context, limit, offset int) ([]*events.Event, error) {
	query := `
		SELECT id, owner_id, title, description, visibility, status,
		       starts_at, ends_at, tz, location, online_url,
		       capacity, waitlist_enabled, created_at, updated_at
		FROM events
		WHERE status = 'published' AND visibility IN ('public', 'by_link')
		ORDER BY starts_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*events.Event
	for rows.Next() {
		var event events.Event
		err := rows.Scan(
			&event.ID, &event.OwnerID, &event.Title, &event.Description,
			&event.Visibility, &event.Status, &event.StartsAt, &event.EndsAt,
			&event.Timezone, &event.Location, &event.OnlineURL,
			&event.Capacity, &event.Waitlist, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &event)
	}

	return result, rows.Err()
}

func (r *EventRepo) ListByUser(ctx context.Context, userID shared.ID) ([]*events.Event, error) {
	query := `
		SELECT e.id, e.owner_id, e.title, e.description, e.visibility, e.status,
		       e.starts_at, e.ends_at, e.tz, e.location, e.online_url,
		       e.capacity, e.waitlist_enabled, e.created_at, e.updated_at
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

	var result []*events.Event
	for rows.Next() {
		var event events.Event
		err := rows.Scan(
			&event.ID, &event.OwnerID, &event.Title, &event.Description,
			&event.Visibility, &event.Status, &event.StartsAt, &event.EndsAt,
			&event.Timezone, &event.Location, &event.OnlineURL,
			&event.Capacity, &event.Waitlist, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &event)
	}

	return result, rows.Err()
}

func (r *EventRepo) Delete(ctx context.Context, id shared.ID) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := r.db.pool.Exec(ctx, query, id)
	return err
}

type RoleRepo struct {
	db *DB
}

func NewRoleRepo(db *DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) Create(ctx context.Context, eventID, userID shared.ID, role events.Role) error {
	query := `INSERT INTO event_roles (id, event_id, user_id, role) VALUES (gen_random_uuid(), $1, $2, $3)`
	_, err := r.db.pool.Exec(ctx, query, eventID, userID, role)
	return err
}

func (r *RoleRepo) GetUserRole(ctx context.Context, eventID, userID shared.ID) (events.Role, error) {
	query := `SELECT role FROM event_roles WHERE event_id = $1 AND user_id = $2`
	var role events.Role
	err := r.db.pool.QueryRow(ctx, query, eventID, userID).Scan(&role)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	return role, err
}

func (r *RoleRepo) ListByEvent(ctx context.Context, eventID shared.ID) (map[shared.ID]events.Role, error) {
	query := `SELECT user_id, role FROM event_roles WHERE event_id = $1`
	rows, err := r.db.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[shared.ID]events.Role)
	for rows.Next() {
		var userID shared.ID
		var role events.Role
		if err := rows.Scan(&userID, &role); err != nil {
			return nil, err
		}
		result[userID] = role
	}

	return result, rows.Err()
}

func (r *RoleRepo) Delete(ctx context.Context, eventID, userID shared.ID) error {
	query := `DELETE FROM event_roles WHERE event_id = $1 AND user_id = $2`
	_, err := r.db.pool.Exec(ctx, query, eventID, userID)
	return err
}

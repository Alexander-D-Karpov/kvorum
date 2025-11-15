package repo

import (
	"context"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/forms"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type FormRepo struct {
	db *DB
}

func NewFormRepo(db *DB) *FormRepo {
	return &FormRepo{db: db}
}

func (r *FormRepo) Create(ctx context.Context, form *forms.Form) error {
	query := `
		INSERT INTO forms (id, event_id, version, schema, rules, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.pool.Exec(ctx, query,
		form.ID, form.EventID, form.Version, form.Schema, form.Rules,
		form.Active, form.CreatedAt, form.UpdatedAt,
	)
	return err
}

func (r *FormRepo) GetByID(ctx context.Context, id shared.ID) (*forms.Form, error) {
	query := `
		SELECT id, event_id, version, schema, rules, active, created_at, updated_at
		FROM forms
		WHERE id = $1
	`

	var form forms.Form
	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&form.ID, &form.EventID, &form.Version, &form.Schema, &form.Rules,
		&form.Active, &form.CreatedAt, &form.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, forms.ErrFormNotFound
	}
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *FormRepo) GetActiveByEvent(ctx context.Context, eventID shared.ID) (*forms.Form, error) {
	query := `
		SELECT id, event_id, version, schema, rules, active, created_at, updated_at
		FROM forms
		WHERE event_id = $1 AND active = true
		ORDER BY created_at DESC
		LIMIT 1
	`

	var form forms.Form
	err := r.db.pool.QueryRow(ctx, query, eventID).Scan(
		&form.ID, &form.EventID, &form.Version, &form.Schema, &form.Rules,
		&form.Active, &form.CreatedAt, &form.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, forms.ErrFormNotFound
	}
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *FormRepo) Update(ctx context.Context, form *forms.Form) error {
	query := `
		UPDATE forms
		SET schema = $2, rules = $3, active = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.pool.Exec(ctx, query,
		form.ID, form.Schema, form.Rules, form.Active, form.UpdatedAt,
	)
	return err
}

type ResponseRepo struct {
	db *DB
}

func NewResponseRepo(db *DB) *ResponseRepo {
	return &ResponseRepo{db: db}
}

func (r *ResponseRepo) Create(ctx context.Context, response *forms.Response) error {
	query := `
		INSERT INTO form_responses (id, form_id, user_id, status, answers, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.pool.Exec(ctx, query,
		response.ID, response.FormID, response.UserID, response.Status,
		response.Answers, response.CreatedAt, response.UpdatedAt,
	)
	return err
}

func (r *ResponseRepo) GetByID(ctx context.Context, id shared.ID) (*forms.Response, error) {
	query := `
		SELECT id, form_id, user_id, status, answers, created_at, updated_at
		FROM form_responses
		WHERE id = $1
	`

	var response forms.Response
	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&response.ID, &response.FormID, &response.UserID, &response.Status,
		&response.Answers, &response.CreatedAt, &response.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, forms.ErrResponseNotFound
	}
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *ResponseRepo) GetByFormAndUser(ctx context.Context, formID, userID shared.ID) (*forms.Response, error) {
	query := `
		SELECT id, form_id, user_id, status, answers, created_at, updated_at
		FROM form_responses
		WHERE form_id = $1 AND user_id = $2
	`

	var response forms.Response
	err := r.db.pool.QueryRow(ctx, query, formID, userID).Scan(
		&response.ID, &response.FormID, &response.UserID, &response.Status,
		&response.Answers, &response.CreatedAt, &response.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, forms.ErrResponseNotFound
	}
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *ResponseRepo) Update(ctx context.Context, response *forms.Response) error {
	query := `
		UPDATE form_responses
		SET status = $2, answers = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.pool.Exec(ctx, query,
		response.ID, response.Status, response.Answers, response.UpdatedAt,
	)
	return err
}

func (r *ResponseRepo) ListByForm(ctx context.Context, formID shared.ID) ([]*forms.Response, error) {
	query := `
		SELECT id, form_id, user_id, status, answers, created_at, updated_at
		FROM form_responses
		WHERE form_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.pool.Query(ctx, query, formID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*forms.Response
	for rows.Next() {
		var response forms.Response
		err := rows.Scan(
			&response.ID, &response.FormID, &response.UserID, &response.Status,
			&response.Answers, &response.CreatedAt, &response.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &response)
	}

	return result, rows.Err()
}

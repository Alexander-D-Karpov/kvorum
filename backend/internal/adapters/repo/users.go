package repo

import (
	"context"

	"github.com/Alexander-D-Karpov/kvorum/internal/app/identity"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *identity.User) error {
	query := `
		INSERT INTO users (id, created_at, updated_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.pool.Exec(ctx, query, user.ID, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	identityQuery := `
		INSERT INTO user_identities (id, user_id, provider, provider_user_id)
		VALUES (gen_random_uuid(), $1, $2, $3)
	`
	_, err = r.db.pool.Exec(ctx, identityQuery, user.ID, user.Provider, user.ProviderID)
	if err != nil {
		return err
	}

	profileQuery := `
		INSERT INTO user_profiles (user_id, display_name, email, phone, tz, locale)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = r.db.pool.Exec(ctx, profileQuery, user.ID, user.DisplayName, user.Email, user.Phone, user.Timezone, user.Locale)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id shared.ID) (*identity.User, error) {
	query := `
		SELECT u.id, u.created_at, u.updated_at, 
		       ui.provider, ui.provider_user_id,
		       up.display_name, up.email, up.phone, up.tz, up.locale
		FROM users u
		JOIN user_identities ui ON ui.user_id = u.id
		LEFT JOIN user_profiles up ON up.user_id = u.id
		WHERE u.id = $1
		LIMIT 1
	`

	var user identity.User
	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
		&user.Provider, &user.ProviderID,
		&user.DisplayName, &user.Email, &user.Phone, &user.Timezone, &user.Locale,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetByProvider(ctx context.Context, provider, providerID string) (*identity.User, error) {
	query := `
		SELECT u.id, u.created_at, u.updated_at,
		       ui.provider, ui.provider_user_id,
		       up.display_name, up.email, up.phone, up.tz, up.locale
		FROM users u
		JOIN user_identities ui ON ui.user_id = u.id
		LEFT JOIN user_profiles up ON up.user_id = u.id
		WHERE ui.provider = $1 AND ui.provider_user_id = $2
		LIMIT 1
	`

	var user identity.User
	err := r.db.pool.QueryRow(ctx, query, provider, providerID).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
		&user.Provider, &user.ProviderID,
		&user.DisplayName, &user.Email, &user.Phone, &user.Timezone, &user.Locale,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *identity.User) error {
	query := `
		UPDATE user_profiles
		SET display_name = $2, email = $3, phone = $4, tz = $5, locale = $6
		WHERE user_id = $1
	`
	_, err := r.db.pool.Exec(ctx, query, user.ID, user.DisplayName, user.Email, user.Phone, user.Timezone, user.Locale)
	return err
}

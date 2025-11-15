package repo

import (
	"context"

	"github.com/Alexander-D-Karpov/kvorum/internal/app/campaigns"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type CampaignRepo struct {
	db *DB
}

func NewCampaignRepo(db *DB) *CampaignRepo {
	return &CampaignRepo{db: db}
}

func (r *CampaignRepo) Create(ctx context.Context, campaign *campaigns.Campaign) error {
	query := `
        INSERT INTO campaigns (id, event_id, name, segment, content, schedule_at, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := r.db.pool.Exec(ctx, query,
		campaign.ID,
		campaign.EventID,
		campaign.Name,
		campaign.Segment,
		campaign.Content,
		campaign.ScheduleAt,
		campaign.Status,
		campaign.CreatedAt,
		campaign.UpdatedAt,
	)
	return err
}

func (r *CampaignRepo) GetByID(ctx context.Context, id shared.ID) (*campaigns.Campaign, error) {
	query := `
        SELECT id, event_id, name, segment, content, schedule_at, status, created_at, updated_at
        FROM campaigns
        WHERE id = $1
    `

	var campaign campaigns.Campaign
	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&campaign.ID,
		&campaign.EventID,
		&campaign.Name,
		&campaign.Segment,
		&campaign.Content,
		&campaign.ScheduleAt,
		&campaign.Status,
		&campaign.CreatedAt,
		&campaign.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

func (r *CampaignRepo) ListByEvent(ctx context.Context, eventID shared.ID) ([]*campaigns.Campaign, error) {
	query := `
        SELECT id, event_id, name, segment, content, schedule_at, status, created_at, updated_at
        FROM campaigns
        WHERE event_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*campaigns.Campaign
	for rows.Next() {
		var campaign campaigns.Campaign
		err := rows.Scan(
			&campaign.ID,
			&campaign.EventID,
			&campaign.Name,
			&campaign.Segment,
			&campaign.Content,
			&campaign.ScheduleAt,
			&campaign.Status,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &campaign)
	}

	return result, rows.Err()
}

func (r *CampaignRepo) Update(ctx context.Context, campaign *campaigns.Campaign) error {
	query := `
        UPDATE campaigns
        SET status = $2, updated_at = $3
        WHERE id = $1
    `
	_, err := r.db.pool.Exec(ctx, query, campaign.ID, campaign.Status, campaign.UpdatedAt)
	return err
}

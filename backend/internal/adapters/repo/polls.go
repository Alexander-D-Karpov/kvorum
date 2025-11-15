package repo

import (
	"context"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/polls"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type PollRepo struct {
	db *DB
}

func NewPollRepo(db *DB) *PollRepo {
	return &PollRepo{db: db}
}

func (r *PollRepo) Create(ctx context.Context, poll *polls.Poll) error {
	query := `
		INSERT INTO polls (id, event_id, question, options, type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.pool.Exec(ctx, query,
		poll.ID, poll.EventID, poll.Question, poll.Options, poll.Type, poll.CreatedAt,
	)
	return err
}

func (r *PollRepo) GetByID(ctx context.Context, id shared.ID) (*polls.Poll, error) {
	query := `
		SELECT id, event_id, question, options, type, created_at, updated_at
		FROM polls
		WHERE id = $1
	`

	var poll polls.Poll
	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&poll.ID, &poll.EventID, &poll.Question, &poll.Options, &poll.Type,
		&poll.CreatedAt, &poll.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, polls.ErrPollNotFound
	}
	if err != nil {
		return nil, err
	}

	return &poll, nil
}

func (r *PollRepo) ListByEvent(ctx context.Context, eventID shared.ID) ([]*polls.Poll, error) {
	query := `
		SELECT id, event_id, question, options, type, created_at, updated_at
		FROM polls
		WHERE event_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*polls.Poll
	for rows.Next() {
		var poll polls.Poll
		err := rows.Scan(
			&poll.ID, &poll.EventID, &poll.Question, &poll.Options, &poll.Type,
			&poll.CreatedAt, &poll.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &poll)
	}

	return result, rows.Err()
}

type VoteRepo struct {
	db *DB
}

func NewVoteRepo(db *DB) *VoteRepo {
	return &VoteRepo{db: db}
}

func (r *VoteRepo) Create(ctx context.Context, vote *polls.Vote) error {
	query := `
		INSERT INTO poll_votes (id, poll_id, user_id, option_key, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.pool.Exec(ctx, query,
		vote.ID, vote.PollID, vote.UserID, vote.OptionKey, vote.CreatedAt,
	)
	return err
}

func (r *VoteRepo) GetByPollAndUser(ctx context.Context, pollID, userID shared.ID) (*polls.Vote, error) {
	query := `
		SELECT id, poll_id, user_id, option_key, created_at
		FROM poll_votes
		WHERE poll_id = $1 AND user_id = $2
	`

	var vote polls.Vote
	err := r.db.pool.QueryRow(ctx, query, pollID, userID).Scan(
		&vote.ID, &vote.PollID, &vote.UserID, &vote.OptionKey, &vote.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &vote, nil
}

func (r *VoteRepo) CountByOption(ctx context.Context, pollID shared.ID) (map[string]int, error) {
	query := `
		SELECT option_key, COUNT(*)
		FROM poll_votes
		WHERE poll_id = $1
		GROUP BY option_key
	`

	rows, err := r.db.pool.Query(ctx, query, pollID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var optionKey string
		var count int
		if err := rows.Scan(&optionKey, &count); err != nil {
			return nil, err
		}
		result[optionKey] = count
	}

	return result, rows.Err()
}

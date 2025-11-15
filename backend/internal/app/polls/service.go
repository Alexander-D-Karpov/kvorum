package polls

import (
	"context"
	"encoding/json"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/polls"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type PollRepo interface {
	Create(ctx context.Context, poll *polls.Poll) error
	GetByID(ctx context.Context, id shared.ID) (*polls.Poll, error)
	ListByEvent(ctx context.Context, eventID shared.ID) ([]*polls.Poll, error)
}

type VoteRepo interface {
	Create(ctx context.Context, vote *polls.Vote) error
	GetByPollAndUser(ctx context.Context, pollID, userID shared.ID) (*polls.Vote, error)
	CountByOption(ctx context.Context, pollID shared.ID) (map[string]int, error)
}

type Service struct {
	pollRepo PollRepo
	voteRepo VoteRepo
}

func NewService(pollRepo PollRepo, voteRepo VoteRepo) *Service {
	return &Service{
		pollRepo: pollRepo,
		voteRepo: voteRepo,
	}
}

func (s *Service) CreatePoll(ctx context.Context, eventID shared.ID, question string, options json.RawMessage, pollType interface{}) (interface{}, error) {
	pt := polls.PollTypeSingle
	if pollType != nil {
		pt = pollType.(polls.PollType)
	}

	poll := polls.NewPoll(eventID, question, options, pt)
	if err := s.pollRepo.Create(ctx, poll); err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *Service) Vote(ctx context.Context, pollID, userID shared.ID, optionKey string) error {
	existing, err := s.voteRepo.GetByPollAndUser(ctx, pollID, userID)
	if err != nil {
		return err
	}

	if existing != nil {
		return polls.ErrAlreadyVoted
	}

	vote := polls.NewVote(pollID, userID, optionKey)
	return s.voteRepo.Create(ctx, vote)
}

func (s *Service) GetResults(ctx context.Context, pollID shared.ID) (map[string]int, error) {
	return s.voteRepo.CountByOption(ctx, pollID)
}

func (s *Service) GetPollsByEvent(ctx context.Context, eventID shared.ID) (interface{}, error) {
	return s.pollRepo.ListByEvent(ctx, eventID)
}

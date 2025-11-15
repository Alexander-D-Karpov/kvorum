package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type AsynqScheduler struct {
	client *asynq.Client
}

func NewAsynqScheduler(redisURL string) (*AsynqScheduler, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})

	return &AsynqScheduler{client: client}, nil
}

func (a *AsynqScheduler) ScheduleReminder(ctx context.Context, at time.Time, payload interface{}) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	task := asynq.NewTask("reminder", data)
	info, err := a.client.Enqueue(task, asynq.ProcessAt(at))
	if err != nil {
		return "", err
	}

	return info.ID, nil
}

func (a *AsynqScheduler) ScheduleCampaign(ctx context.Context, campaignID string, at time.Time) error {
	data, _ := json.Marshal(map[string]string{"campaign_id": campaignID})
	task := asynq.NewTask("campaign", data)
	_, err := a.client.Enqueue(task, asynq.ProcessAt(at))
	return err
}

func (a *AsynqScheduler) Close() error {
	return a.client.Close()
}

type AsynqServer struct {
	server *asynq.Server
}

func NewAsynqServer(redisURL string) (*AsynqServer, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     opts.Addr,
			Password: opts.Password,
			DB:       opts.DB,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &AsynqServer{server: srv}, nil
}

func (a *AsynqServer) HandleFunc(pattern string, handler func(context.Context, *asynq.Task) error) {
	mux := asynq.NewServeMux()
	mux.HandleFunc(pattern, handler)
	a.server.Start(mux)
}

func (a *AsynqServer) Start(mux *asynq.ServeMux) error {
	return a.server.Start(mux)
}

func (a *AsynqServer) Stop() {
	a.server.Stop()
	a.server.Shutdown()
}

type TaskHandler struct {
	botClient interface {
		SendMessage(ctx context.Context, chatID int64, req interface{}) (string, error)
	}
}

func (h *TaskHandler) HandleReminder(ctx context.Context, task *asynq.Task) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	return nil
}

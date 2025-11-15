package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (r *RedisCache) GetEventPublic(ctx context.Context, id shared.ID) (*events.Event, bool) {
	key := fmt.Sprintf("event:pub:%s", id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	var event events.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, false
	}

	return &event, true
}

func (r *RedisCache) SetEventPublic(ctx context.Context, id shared.ID, event *events.Event, ttl time.Duration) {
	key := fmt.Sprintf("event:pub:%s", id)
	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	r.client.Set(ctx, key, data, ttl)
}

func (r *RedisCache) InvalidateEvent(ctx context.Context, id shared.ID) {
	key := fmt.Sprintf("event:pub:%s", id)
	r.client.Del(ctx, key)
}

func (r *RedisCache) GetDraft(ctx context.Context, formID, userID shared.ID) (json.RawMessage, bool) {
	key := fmt.Sprintf("draft:form:%s:%s", formID, userID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return data, true
}

func (r *RedisCache) SetDraft(ctx context.Context, formID, userID shared.ID, data json.RawMessage, ttl time.Duration) {
	key := fmt.Sprintf("draft:form:%s:%s", formID, userID)
	r.client.Set(ctx, key, data, ttl)
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

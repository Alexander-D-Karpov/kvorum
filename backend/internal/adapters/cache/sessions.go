package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/security"
)

func (r *RedisCache) SetSession(ctx context.Context, session *security.Session) error {
	key := fmt.Sprintf("session:%s", session.ID)
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	ttl := time.Until(session.ExpiresAt)
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache) GetSession(ctx context.Context, sessionID string) (*security.Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var session security.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	if session.IsExpired() {
		r.DeleteSession(ctx, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return &session, nil
}

func (r *RedisCache) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return r.client.Del(ctx, key).Err()
}

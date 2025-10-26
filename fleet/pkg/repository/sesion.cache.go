package repository

import (
	"context"
	"encoding/json"

	"github.com/elzestia/fleet/pkg/models"
)

func buildSessionCacheKey(userId, sessionId string) string {
	return "session:" + userId + ":" + sessionId
}

func (r *Repository) SetSessionCache(ctx context.Context, session models.Session) (err error) {
	key := buildSessionCacheKey(session.UserID, session.ID)
	jsonData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = r.redisConn.Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetSessionCache(ctx context.Context, userId, sessionId string) (session models.Session, err error) {
	key := buildSessionCacheKey(userId, sessionId)
	data, err := r.redisConn.Get(ctx, key).Bytes()
	if err != nil {
		return session, err
	}

	err = json.Unmarshal(data, &session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func (r *Repository) DeleteSessionCache(ctx context.Context, userId, sessionId string) (err error) {
	key := buildSessionCacheKey(userId, sessionId)
	err = r.redisConn.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

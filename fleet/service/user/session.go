package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"go.uber.org/zap"
)

func (s *UserService) insertSessions(ctx context.Context, session models.InsertSessionParams) (err error) {
	err = s.repository.InsertSession(ctx, session)
	if err != nil {
		s.logger.Error("[insertSessions] failed to insert session", zap.Error(err))
	}

	return
}

func (s *UserService) GetSessionByID(ctx context.Context, req *models.AuthClaims) (session models.Session, err error) {
	if req == nil {
		s.logger.Error("[GetSessionByID] Invalid request for GetSessionByID: req is nil")
		return session, response.ErrorUnAuthorized
	}
	if req.ID == "" || req.Subject == "" {
		s.logger.Error("[GetSessionByID] Invalid request for GetSessionByID", zap.String("ID", req.ID), zap.String("Subject", req.Subject))
		return session, response.ErrorUnAuthorized
	}

	// Check cache first
	session, err = s.repository.GetSessionCache(ctx, req.Subject, req.ID)
	if err == nil {
		s.logger.Debug("[GetSessionByID] Session found in cache", zap.String("ID", req.ID), zap.String("UserID", req.Subject))
		return session, nil
	}

	s.logger.Warn("[GetSessionByID] Failed to get session from cache", zap.Error(err))
	err = nil

	session, err = s.repository.GetSessionByEntityId(ctx, models.GetSessionByEntityIdParams{
		ID:     req.ID,
		UserID: req.Subject,
	})
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.Debug("[GetSessionByID] Session not found", zap.String("ID", req.ID), zap.String("UserID", req.Subject))
		err = response.ErrorUnAuthorized
		return
	}

	if err != nil {
		s.logger.Error("[GetSessionByID] Failed to get session by id", zap.Error(err))
		err = response.ErrorInternalServerError
	}

	go func() {
		// Cache the session
		err = s.repository.SetSessionCache(context.Background(), session)
		if err != nil {
			s.logger.Error("[GetSessionByID] Failed to set session cache", zap.Error(err))
		}
	}()

	return
}

func (s *UserService) deleteSessionByID(ctx context.Context, id, userId string) (err error) {
	err = s.repository.DeleteSessionByID(ctx, models.DeleteSessionByIDParams{
		ID:     id,
		UserID: userId,
	})
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.Debug("[deleteSessionByID] Session not found", zap.Error(err))
		err = response.ErrorUnAuthorized
		return
	}

	if err != nil {
		s.logger.Error("[deleteSessionByID] Failed to delete session by id", zap.Error(err))
		err = response.ErrorInternalServerError
	}

	return
}

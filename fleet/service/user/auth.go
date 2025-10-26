package user

import (
	"context"
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/transport/http/request"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/elzestia/fleet/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *UserService) Login(ctx context.Context, req *request.Login) (res response.Login, err error) {
	user, err := s.GetUserByEmail(ctx, req.Email)
	if errors.Is(err, response.ErrorUserDatabaseUserNotFound) {
		err = response.ErrorInvalidEmailOrPassword
		return
	}

	if err != nil {
		return
	}

	err = s.findAndCompareUserPassword(ctx, user.ID, req.Password)
	if err != nil {
		return
	}

	token, idToken, expiresAt, err := s.generateToken(s.config.Function.Auth.Token.SecretKey, s.config.Function.Auth.Token.Expire, user)
	if err != nil {
		s.logger.Error("[Login] Failed to generate token", zap.Error(err))
		err = response.ErrorInternalServerError
		return
	}

	session := models.InsertSessionParams{
		ID: idToken,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
		UserID:    user.ID,
		CreatedBy: user.ID,
	}

	err = s.insertSessions(ctx, session)
	if err != nil {
		s.logger.Error("[Login] Failed to insert session", zap.Error(err))
		err = response.ErrorInternalServerError
		return
	}

	res.Token = token
	return
}

func (s *UserService) VerifyToken(ctx context.Context, token string) (claims models.AuthClaims, err error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy(s.config.App.Name))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	secretKeyParsed, err := paseto.V4SymmetricKeyFromHex(s.config.Function.Auth.Token.SecretKey)
	if err != nil {
		s.logger.Error("[VerifyToken] Failed to parse public key", zap.Error(err))
		err = response.InvalidToken
		return
	}

	parsedToken, err := parser.ParseV4Local(secretKeyParsed, token, nil)
	if err != nil {
		s.logger.Error("[VerifyToken] Failed to parse token", zap.Error(err))
		err = response.InvalidToken
		return
	}

	rawData := parsedToken.Claims()
	sessionIDString, err := util.DecryptData([]byte(rawData["id"].(string)), []byte(s.config.Function.Auth.SecretKey.SessionID), 30)
	if err != nil {
		s.logger.Error("[VerifyToken] Failed to decrypt session id", zap.Error(err))
		err = response.InvalidToken
		return
	}

	subject, ok := rawData["sub"].(string)
	if !ok {
		s.logger.Error("[VerifyToken] Failed to convert subject to string", zap.Error(err))
		err = response.InvalidToken
		return
	}

	sessionID := sessionIDString
	claims.ID = string(sessionID)
	claims.Subject = subject

	return claims, nil
}

func (s *UserService) Logout(ctx context.Context, req *models.AuthClaims) (err error) {
	if req == nil {
		s.logger.Error("[Logout] Invalid request: req is nil")
		return response.ErrorUnAuthorized
	}

	if req.ID == "" || req.Subject == "" {
		s.logger.Error("[Logout] Invalid request", zap.String("ID", req.ID), zap.String("Subject", req.Subject))
		return response.ErrorUnAuthorized
	}

	err = s.deleteSessionByID(ctx, req.ID, req.Subject)
	if err != nil {
		s.logger.Error("[Logout] Failed to delete session by id", zap.Error(err))
		return response.ErrorInternalServerError
	}

	err = s.deleteSessionByID(ctx, req.Subject, req.ID)
	if err != nil {
		s.logger.Error("[Logout] Failed to delete session from cache", zap.Error(err))
		return response.ErrorInternalServerError
	}

	return nil
}

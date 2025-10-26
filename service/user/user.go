package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/transport/http/request"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/elzestia/fleet/pkg/util"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
)

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	hashedEmail := util.Hash(email, s.config.Function.User.SecretKey.EmailSalt)
	user, err = s.repository.GetUserByEmail(ctx, hashedEmail)
	if errors.Is(err, sql.ErrNoRows) {
		err = response.ErrorUserDatabaseUserNotFound
		return
	}

	if err != nil {
		s.logger.Error("[GetUserByEmail] Failed to get user by params", zap.Error(err))
		err = response.ErrorInternalServerError
		return
	}

	return
}

func (s *UserService) GetUserByUserID(ctx context.Context, userID string) (userDetail response.UserDetail, err error) {
	user, err := s.repository.GetUserByID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		err = response.ErrorUserDatabaseUserNotFound
		return
	}

	if err != nil {
		s.logger.Error("[GetUserByUserID] Failed to get user by ID", zap.Error(err))
		err = response.ErrorInternalServerError
		return
	}

	// Decrypt email
	decryptedEmail, err := util.DecryptData([]byte(user.Email), []byte(s.config.Function.User.SecretKey.Email), s.config.Function.User.SecretKey.EmailSaltLength)
	if err != nil {
		s.logger.Error("[GetUserByUserID] Failed to decrypt email", zap.Error(err))
		err = response.ErrorInternalServerError
		return
	}

	updatedAt := ""
	if user.UpdatedAt.Valid {
		updatedAt = user.UpdatedAt.Time.Format(time.RFC3339)
	}

	userDetail = response.UserDetail{
		ID:        user.ID,
		Email:     string(decryptedEmail),
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: updatedAt,
	}

	return
}

func (s *UserService) RegisterUser(ctx context.Context, req *request.RegisterUserRequest) (err error) {
	encryptedEmail, err := util.EncryptData([]byte(req.Email), []byte(s.config.Function.User.SecretKey.Email), s.config.Function.User.SecretKey.EmailSaltLength)
	if err != nil {
		s.logger.Error("[RegisterUser] Failed to encrypt email", zap.Error(err))
		return response.ErrorInternalServerError
	}

	user, err := s.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, response.ErrorUserDatabaseUserNotFound) {
		s.logger.Error("[RegisterUser] Failed to get user by params", zap.Error(err))
		return response.ErrorInternalServerError
	}
	if user.ID != "" {
		return response.ErrorUserDatabaseUserEmailAlreadyUsed
	}

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		s.logger.Error("[RegisterUser] Failed to begin transaction", zap.Error(err))
		return response.ErrorInternalServerError
	}

	defer func() {
		if err != nil {
			if rollbackErr := s.repository.RollbackTx(tx); rollbackErr != nil {
				s.logger.Error("[RegisterUser] Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	id := ulid.Make()
	hashedEmail := util.Hash(req.Email, s.config.Function.User.SecretKey.EmailSalt)
	txQueries := s.repository.WithTx(tx)
	err = txQueries.InsertUser(ctx, models.InsertUserParams{
		ID:        id.String(),
		Email:     encryptedEmail,
		EmailHash: hashedEmail,
		Name:      req.Name,
		CreatedBy: id.String(),
	})
	if err != nil {
		s.logger.Error("[RegisterUser] Failed to insert user", zap.Error(err))
		return response.ErrorInternalServerError
	}

	passwordHash, err := util.HashPassword(req.Password, s.config.Function.Auth.SecretKey.PasswordSalt)
	if err != nil {
		s.logger.Error("[RegisterUser]Failed to hash password", zap.Error(err))
		return response.ErrorInternalServerError
	}

	authId := ulid.Make()
	err = txQueries.InsertAuth(ctx, models.InsertAuthParams{
		ID:        authId.String(),
		UserID:    id.String(),
		Password:  string(passwordHash),
		CreatedBy: id.String(),
	})
	if err != nil {
		s.logger.Error("[RegisterUser] Failed to insert auth", zap.Error(err))
		return response.ErrorInternalServerError
	}

	err = s.repository.CommitTx(tx)
	if err != nil {
		s.logger.Error("[RegisterUser] Failed to commit transaction", zap.Error(err))
		return response.ErrorInternalServerError
	}

	return nil
}

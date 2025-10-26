package user

import (
	"context"
	"time"

	"aidanwoods.dev/go-paseto"
	"go.uber.org/zap"

	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/elzestia/fleet/pkg/util"
	"github.com/oklog/ulid/v2"
)

func (s *UserService) generateToken(secretKey string, expiration time.Duration, user models.User) (string, string, time.Time, error) {
	token := paseto.NewToken()
	idToken, err := s.buildTokenMetadata(token, expiration, user)
	if err != nil {
		return "", "", time.Time{}, err
	}

	secretKeyParsed, err := paseto.V4SymmetricKeyFromHex(secretKey)
	if err != nil {
		s.logger.Error("[generateToken] Failed to parse secret key", zap.Error(err))
		return "", "", time.Time{}, err
	}

	expiredAt, err := token.GetExpiration()
	if err != nil || expiredAt.Before(time.Now()) || expiration <= 0 {
		s.logger.Error("[generateToken] Failed to get expiration", zap.Error(err))
		err = response.ErrorInternalServerError
		return "", "", time.Time{}, err
	}

	return token.V4Encrypt(secretKeyParsed, nil), idToken, expiredAt, err
}

func (s *UserService) buildTokenMetadata(token paseto.Token, expiration time.Duration, user models.User) (string, error) {
	now := time.Now()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(expiration))
	token.SetIssuer(s.config.App.Name)

	id := ulid.Make()
	idString, err := util.EncryptData([]byte(id.String()), []byte(s.config.Function.Auth.SecretKey.SessionID), 30)
	if err != nil {
		return "", err
	}

	token.SetSubject(user.ID)
	token.SetString("id", idString)

	return id.String(), err
}

func (s *UserService) findAndCompareUserPassword(ctx context.Context, userID string, password string) (err error) {
	auth, err := s.repository.GetAuthByUserUserID(ctx, userID)
	if err != nil {
		s.logger.Error("[findAndCompareUserPassword] Failed to get auth by user user id", zap.Error(err))
		err = response.ErrorInternalServerError
		return
	}

	err = util.ComparePassword(password, s.config.Function.Auth.SecretKey.PasswordSalt, auth.Password)
	if err != nil {
		err = response.ErrorInvalidEmailOrPassword
	}

	return
}

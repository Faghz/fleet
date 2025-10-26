package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/mocks"
	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/transport/http/request"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/elzestia/fleet/pkg/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Table-driven test for findAndCompareUserPassword
func TestUserService_findAndCompareUserPassword(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		password      string
		mockReturn    models.Auth
		mockError     error
		expectedError error
	}{
		{
			name:          "Auth not found",
			userID:        "user-123",
			password:      "password123",
			mockReturn:    models.Auth{},
			mockError:     errors.New("auth not found"),
			expectedError: response.ErrorInternalServerError,
		},
		{
			name:     "Password mismatch",
			userID:   "user-123",
			password: "wrongpassword",
			mockReturn: models.Auth{
				ID:       "auth-123",
				UserID:   "user-123",
				Password: "$2a$10$hashedpassword",
			},
			mockError:     nil,
			expectedError: nil, // Will be set by password comparison logic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _ := createTestUserService(t)
			ctx := context.Background()

			mockRepo.EXPECT().GetAuthByUserUserID(ctx, tt.userID).Return(tt.mockReturn, tt.mockError)

			err := service.findAndCompareUserPassword(ctx, tt.userID, tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				// For password mismatch, we expect an error but don't check the specific type
				// since it depends on the password comparison logic
				assert.Error(t, err)
			}
		})
	}
}

// Table-driven test for Login
func TestUserService_Login(t *testing.T) {
	defaultConfig := &configs.Config{
		Function: configs.FunctionConfig{
			Auth: configs.FunctionAuth{
				SecretKey: configs.FunctionAuthSecretKey{
					PasswordSalt: "ohirO31p28iP",
					SessionID:    "1LytH6biNjlNmNOXb9Yud7ng3hlAKnCG",
				},
				Token: configs.FunctionAuthToken{
					SecretKey: "25e1a26721010a409c6a76ca80cba805ffedc32fb0c7cd3b1ed34cdfd283ba97",
					Expire:    24 * time.Hour,
				},
			},
		},
	}

	type testCase struct {
		name      string
		loginReq  *request.Login
		config    *configs.Config
		setupAuth bool
		calls     func(context.Context, *testing.T, *mocks.MockUserRepository, testCase)
		expect    func(context.Context, *testing.T, testCase, response.Login, error)
	}
	tests := []testCase{
		{
			name: "Invalid email - user not found",
			loginReq: &request.Login{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},

			setupAuth: false,
			calls: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, tt testCase) {
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, sql.ErrNoRows)
			},
			expect: func(ctx context.Context, t *testing.T, tt testCase, res response.Login, err error) {
				assert.Empty(t, res.Token)
				assert.Error(t, err)
				assert.Equal(t, response.ErrorInvalidEmailOrPassword, err)
			},
		},
		{
			name: "Database error on user fetch",
			loginReq: &request.Login{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupAuth: false,
			calls: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, tt testCase) {
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, errors.New("database connection error"))
			},
			expect: func(ctx context.Context, t *testing.T, tt testCase, res response.Login, err error) {
				assert.Empty(t, res.Token)
				assert.Error(t, err)
				assert.Equal(t, response.ErrorInternalServerError, err)
			},
		},
		{
			name: "Auth fetch error after user found",
			loginReq: &request.Login{
				Email:    "test@example.com",
				Password: "password123",
			},

			setupAuth: true,
			calls: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, tt testCase) {
				user := models.User{
					ID:    "user-123",
					Email: "test@example.com",
				}
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil)
				mockRepo.EXPECT().GetAuthByUserUserID(ctx, user.ID).Return(models.Auth{
					ID:       "auth-123",
					UserID:   "user-123",
					Password: "$2a$10$hashedpassword", // Assume this is a hashed
				}, errors.New("auth fetch failed"))
			},
			expect: func(ctx context.Context, t *testing.T, tt testCase, res response.Login, err error) {
				assert.Empty(t, res.Token)
				assert.Error(t, err)
				assert.Equal(t, response.ErrorInternalServerError, err)
			},
		},
		{
			name: "Password mismatch",
			loginReq: &request.Login{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupAuth: true,
			calls: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, tt testCase) {
				user := models.User{
					ID:    "user-123",
					Email: "test@example.com",
				}
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil)
				mockRepo.EXPECT().GetAuthByUserUserID(ctx, user.ID).Return(models.Auth{
					ID:       "auth-123",
					UserID:   "user-123",
					Password: "$2a$10$hashedpassword", // Assume this is a hashed
				}, nil)
			},
			expect: func(ctx context.Context, t *testing.T, tt testCase, res response.Login, err error) {
				assert.Empty(t, res.Token)
				assert.Error(t, err)
				assert.Equal(t, response.ErrorInvalidEmailOrPassword, err)
			},
		},
		{
			name: "Fail encrypt session id",
			loginReq: &request.Login{
				Email:    "test@example.com",
				Password: "IsThisAPassword?12345#@!",
			},
			config: &configs.Config{
				Function: configs.FunctionConfig{
					Auth: configs.FunctionAuth{
						SecretKey: configs.FunctionAuthSecretKey{
							PasswordSalt: "ohirO31p28iP",
						},
						Token: configs.FunctionAuthToken{
							SecretKey: "25e1a26721010a409c6a76ca80cba805ffedc32fb0c7cd3b1ed34cdfd283ba97",
							Expire:    24 * time.Hour,
						},
					},
				},
			},
			setupAuth: true,
			calls: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, tt testCase) {
				user := models.User{
					ID:    "user-123",
					Email: "IsThisAPassword?12345#@!",
				}
				mockRepo.EXPECT().GetAuthByUserUserID(ctx, user.ID).Return(models.Auth{
					ID:       "auth-123",
					UserID:   "user-123",
					Password: "$2a$13$mXH2I2BJo6OYiOsv.FUgSeomg2jnXP1Of7y.9H9WwdMx4QOhImXAq",
				}, nil)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil)
			},
			expect: func(ctx context.Context, t *testing.T, tt testCase, res response.Login, err error) {
				assert.Empty(t, res.Token)
				assert.Error(t, err)
			},
		},

		{
			name: "Successful login",
			loginReq: &request.Login{
				Email:    "test@example.com",
				Password: "IsThisAPassword?12345#@!",
			},
			setupAuth: true,
			calls: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, tt testCase) {
				user := models.User{
					ID:    "user-123",
					Email: "IsThisAPassword?12345#@!",
				}
				mockRepo.EXPECT().GetAuthByUserUserID(ctx, user.ID).Return(models.Auth{
					ID:       "auth-123",
					UserID:   "user-123",
					Password: "$2a$13$mXH2I2BJo6OYiOsv.FUgSeomg2jnXP1Of7y.9H9WwdMx4QOhImXAq",
				}, nil)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil)
				mockRepo.EXPECT().InsertSession(ctx, gomock.Any()).Return(nil)
			},
			expect: func(ctx context.Context, t *testing.T, tt testCase, res response.Login, err error) {
				assert.NotEmpty(t, res.Token)
				assert.NoError(t, err)
			},
		},
		{
			name: "Failed insert session",
			loginReq: &request.Login{
				Email:    "test@example.com",
				Password: "IsThisAPassword?12345#@!",
			},
			setupAuth: true,
			calls: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, tt testCase) {
				user := models.User{
					ID:    "user-123",
					Email: "IsThisAPassword?12345#@!",
				}
				mockRepo.EXPECT().GetAuthByUserUserID(ctx, user.ID).Return(models.Auth{
					ID:       "auth-123",
					UserID:   "user-123",
					Password: "$2a$13$mXH2I2BJo6OYiOsv.FUgSeomg2jnXP1Of7y.9H9WwdMx4QOhImXAq",
				}, nil)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil)
				mockRepo.EXPECT().InsertSession(ctx, gomock.Any()).Return(errors.New("failed to insert session"))
			},
			expect: func(ctx context.Context, t *testing.T, tt testCase, res response.Login, err error) {
				assert.Empty(t, res.Token)
				assert.Equal(t, response.ErrorInternalServerError, err)
			},
		},
		{
			name: "Failed to parse secret key",
			loginReq: &request.Login{
				Email:    "test@example.com",
				Password: "IsThisAPassword?12345#@!",
			},
			config: &configs.Config{
				Function: configs.FunctionConfig{
					Auth: configs.FunctionAuth{
						SecretKey: configs.FunctionAuthSecretKey{
							PasswordSalt: "ohirO31p28iP",
							SessionID:    "1LytH6biNjlNmNOXb9Yud7ng3hlAKnCG",
						},
						Token: configs.FunctionAuthToken{
							SecretKey: "invalid-secret-key", // Invalid key
							Expire:    24 * time.Hour,
						},
					},
				},
			},
			setupAuth: true,
			calls: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, tt testCase) {
				user := models.User{
					ID:    "user-123",
					Email: "IsThisAPassword?12345#@!",
				}
				mockRepo.EXPECT().GetAuthByUserUserID(ctx, user.ID).Return(models.Auth{
					ID:       "auth-123",
					UserID:   "user-123",
					Password: "$2a$13$mXH2I2BJo6OYiOsv.FUgSeomg2jnXP1Of7y.9H9WwdMx4QOhImXAq",
				}, nil)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil)
			},
			expect: func(ctx context.Context, t *testing.T, tt testCase, res response.Login, err error) {
				assert.Empty(t, res.Token)
				assert.Equal(t, response.ErrorInternalServerError, err)
			},
		},
		{
			name: "Failed - Get token expiration",
			loginReq: &request.Login{
				Email:    "test@example.com",
				Password: "IsThisAPassword?12345#@!",
			},
			config: &configs.Config{
				App: configs.AppConfig{
					Name: "test-app",
				},
				Function: configs.FunctionConfig{
					Auth: configs.FunctionAuth{
						SecretKey: configs.FunctionAuthSecretKey{
							PasswordSalt: "ohirO31p28iP",
							SessionID:    "invalid", // Too short for encryption - should cause EncryptData to fail
						},
						Token: configs.FunctionAuthToken{
							SecretKey: "25e1a26721010a409c6a76ca80cba805ffedc32fb0c7cd3b1ed34cdfd283ba97",
							Expire:    24 * time.Hour,
						},
					},
				},
			},
			setupAuth: true,
			calls: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, tt testCase) {
				user := models.User{
					ID:    "user-123",
					Email: "IsThisAPassword?12345#@!",
				}
				mockRepo.EXPECT().GetAuthByUserUserID(ctx, user.ID).Return(models.Auth{
					ID:       "auth-123",
					UserID:   "user-123",
					Password: "$2a$13$mXH2I2BJo6OYiOsv.FUgSeomg2jnXP1Of7y.9H9WwdMx4QOhImXAq",
				}, nil)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil)
			},
			expect: func(ctx context.Context, t *testing.T, tt testCase, res response.Login, err error) {
				assert.Empty(t, res.Token)
				assert.Error(t, err)
				// This should fail during token generation in buildTokenMetadata when encrypting session ID
				assert.Equal(t, response.ErrorInternalServerError, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			service, mockRepo, _ := createTestUserService(t)
			service.config = tt.config
			if tt.config == nil {
				service.config = defaultConfig
			}

			tt.calls(ctx, t, mockRepo, tt)
			res, err := service.Login(ctx, tt.loginReq)
			tt.expect(ctx, t, tt, res, err)
		})
	}
}

// Table-driven test for VerifyToken
func TestUserService_VerifyToken(t *testing.T) {
	// Test configuration with valid keys
	testConfig := &configs.Config{
		App: configs.AppConfig{
			Name: "test-app",
		},
		Function: configs.FunctionConfig{
			Auth: configs.FunctionAuth{
				Token: configs.FunctionAuthToken{
					SecretKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", // 64 hex chars = 32 bytes
					Expire:    time.Hour * 3,                                                      // 1 hour
				},
				SecretKey: configs.FunctionAuthSecretKey{
					SessionID: "0123456789abcdef0123456789abcdef", // 32 bytes for session encryption
				},
			},
		},
	}

	// Helper function to generate a valid token for testing
	generateTestToken := func(config *configs.Config, userID, sessionID string, expiresAt time.Time) (string, error) {
		token := paseto.NewToken()
		token.SetIssuedAt(time.Now())
		token.SetNotBefore(time.Now())
		token.SetExpiration(expiresAt)
		token.SetIssuer(config.App.Name)
		token.SetSubject(userID)

		// Encrypt session ID
		encryptedSessionID, err := util.EncryptData([]byte(sessionID), []byte(config.Function.Auth.SecretKey.SessionID), 30)
		if err != nil {
			return "", err
		}

		token.Set("id", string(encryptedSessionID))

		secretKey, err := paseto.V4SymmetricKeyFromHex(config.Function.Auth.Token.SecretKey)
		if err != nil {
			return "", err
		}

		return token.V4Encrypt(secretKey, nil), nil
	}

	tests := []struct {
		name           string
		setupToken     func() string
		config         *configs.Config
		expectedError  error
		expectedClaims models.AuthClaims
	}{
		{
			name: "Valid token with correct claims",
			setupToken: func() string {
				token, err := generateTestToken(testConfig, "user-123", "session-456", time.Now().Add(time.Hour))
				assert.NoError(t, err)
				return token
			},
			config: testConfig,
			expectedClaims: models.AuthClaims{
				ID:      "session-456",
				Subject: "user-123",
			},
		},
		{
			name: "Expired token",
			setupToken: func() string {
				token, err := generateTestToken(testConfig, "user-123", "session-456", time.Now().Add(-time.Hour))
				assert.NoError(t, err)
				return token
			},
			config:        testConfig,
			expectedError: response.InvalidToken,
		},
		{
			name: "Token with invalid issuer",
			setupToken: func() string {
				invalidConfig := *testConfig
				invalidConfig.App.Name = "wrong-app"
				token, err := generateTestToken(&invalidConfig, "user-123", "session-456", time.Now().Add(time.Hour))
				assert.NoError(t, err)
				return token
			},
			config:        testConfig,
			expectedError: response.InvalidToken,
		},
		{
			name: "Invalid secret key in config",
			setupToken: func() string {
				token, err := generateTestToken(testConfig, "user-123", "session-456", time.Now().Add(time.Hour))
				assert.NoError(t, err)
				return token
			},
			config: &configs.Config{
				App: configs.AppConfig{
					Name: "test-app",
				},
				Function: configs.FunctionConfig{
					Auth: configs.FunctionAuth{
						Token: configs.FunctionAuthToken{
							SecretKey: "invalid-key", // Invalid hex key
						},
						SecretKey: configs.FunctionAuthSecretKey{
							SessionID: "0123456789abcdef",
						},
					},
				},
			},
			expectedError: response.InvalidToken,
		},
		{
			name: "Empty token",
			setupToken: func() string {
				return ""
			},
			config:        testConfig,
			expectedError: response.InvalidToken,
		},
		{
			name: "Malformed token",
			setupToken: func() string {
				return "invalid.token.format"
			},
			config:        testConfig,
			expectedError: response.InvalidToken,
		},
		{
			name: "Token signed with different key",
			setupToken: func() string {
				differentConfig := *testConfig
				differentConfig.Function.Auth.Token.SecretKey = "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
				token, err := generateTestToken(&differentConfig, "user-123", "session-456", time.Now().Add(time.Hour))
				assert.NoError(t, err)
				return token
			},
			config:        testConfig,
			expectedError: response.InvalidToken,
		},
		{
			name: "Token with missing subject",
			setupToken: func() string {
				token := paseto.NewToken()
				token.SetIssuedAt(time.Now())
				token.SetNotBefore(time.Now())
				token.SetExpiration(time.Now().Add(time.Hour))
				token.SetIssuer(testConfig.App.Name)
				// No subject set

				encryptedSessionID, err := util.EncryptData([]byte("session-456"), []byte(testConfig.Function.Auth.SecretKey.SessionID), 30)
				assert.NoError(t, err)
				token.Set("id", string(encryptedSessionID))

				secretKey, err := paseto.V4SymmetricKeyFromHex(testConfig.Function.Auth.Token.SecretKey)
				assert.NoError(t, err)

				return token.V4Encrypt(secretKey, nil)
			},
			config:        testConfig,
			expectedError: response.InvalidToken,
		},
		{
			name: "Token with corrupted session ID encryption",
			setupToken: func() string {
				token := paseto.NewToken()
				token.SetIssuedAt(time.Now())
				token.SetNotBefore(time.Now())
				token.SetExpiration(time.Now().Add(time.Hour))
				token.SetIssuer(testConfig.App.Name)
				token.SetSubject("user-123")
				token.Set("id", "corrupted-encrypted-data")

				secretKey, err := paseto.V4SymmetricKeyFromHex(testConfig.Function.Auth.Token.SecretKey)
				assert.NoError(t, err)

				return token.V4Encrypt(secretKey, nil)
			},
			config:        testConfig,
			expectedError: response.InvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service with test config
			service := &UserService{
				config: tt.config,
				logger: zap.NewNop(),
			}

			ctx := context.Background()
			token := tt.setupToken()

			// Call the function under test
			claims, err := service.VerifyToken(ctx, token)
			// Assertions
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				// Claims should be empty on error
				assert.Empty(t, claims.ID)
				assert.Empty(t, claims.Subject)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedClaims.ID, claims.ID)
				assert.Equal(t, tt.expectedClaims.Subject, claims.Subject)
			}
		})
	}
}

// Test concurrent token verification
func TestUserService_VerifyToken_Concurrent(t *testing.T) {
	testConfig := &configs.Config{
		App: configs.AppConfig{
			Name: "test-app",
		},
		Function: configs.FunctionConfig{
			Auth: configs.FunctionAuth{
				Token: configs.FunctionAuthToken{
					SecretKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				},
				SecretKey: configs.FunctionAuthSecretKey{
					SessionID: "0123456789abcdef0123456789abcdef",
				},
			},
		},
	}

	service := &UserService{
		config: testConfig,
		logger: zap.NewNop(),
	}

	// Generate a valid token
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(time.Hour))
	token.SetIssuer(testConfig.App.Name)
	token.SetSubject("concurrent-user")

	encryptedSessionID, err := util.EncryptData([]byte("concurrent-session"), []byte(testConfig.Function.Auth.SecretKey.SessionID), 30)
	assert.NoError(t, err)
	token.Set("id", string(encryptedSessionID))

	secretKey, err := paseto.V4SymmetricKeyFromHex(testConfig.Function.Auth.Token.SecretKey)
	assert.NoError(t, err)

	tokenString := token.V4Encrypt(secretKey, nil)

	// Run multiple goroutines to test concurrent access
	const numGoroutines = 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			ctx := context.Background()
			claims, err := service.VerifyToken(ctx, tokenString)
			if err != nil {
				results <- err
				return
			}
			if claims.ID != "concurrent-session" || claims.Subject != "concurrent-user" {
				results <- assert.AnError
				return
			}
			results <- nil
		}()
	}

	// Check all results
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}

// Benchmark token verification
func BenchmarkUserService_VerifyToken(b *testing.B) {
	testConfig := &configs.Config{
		App: configs.AppConfig{
			Name: "benchmark-app",
		},
		Function: configs.FunctionConfig{
			Auth: configs.FunctionAuth{
				Token: configs.FunctionAuthToken{
					SecretKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				},
				SecretKey: configs.FunctionAuthSecretKey{
					SessionID: "0123456789abcdef0123456789abcdef",
				},
			},
		},
	}

	service := &UserService{
		config: testConfig,
		logger: zap.NewNop(),
	}

	// Generate a valid token for benchmarking
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(time.Hour))
	token.SetIssuer(testConfig.App.Name)
	token.SetSubject("benchmark-user")

	encryptedSessionID, _ := util.EncryptData([]byte("benchmark-session"), []byte(testConfig.Function.Auth.SecretKey.SessionID), 30)
	token.Set("id", string(encryptedSessionID))

	secretKey, _ := paseto.V4SymmetricKeyFromHex(testConfig.Function.Auth.Token.SecretKey)
	tokenString := token.V4Encrypt(secretKey, nil)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.VerifyToken(ctx, tokenString)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestUserService_Logout(t *testing.T) {
	type testCase struct {
		name      string
		req       *models.AuthClaims
		mockCalls func(*testing.T, *mocks.MockUserRepository)
		wantErr   error
	}

	tests := []testCase{
		{
			name: "Successful logout",
			req: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			mockCalls: func(t *testing.T, mockRepo *mocks.MockUserRepository) {
				// First call: delete session by ID
				mockRepo.EXPECT().
					DeleteSessionByID(gomock.Any(), models.DeleteSessionByIDParams{
						ID:     "session-123",
						UserID: "user-456",
					}).
					Return(nil)

				// Second call: delete session from cache (swapped parameters)
				mockRepo.EXPECT().
					DeleteSessionByID(gomock.Any(), models.DeleteSessionByIDParams{
						ID:     "user-456",
						UserID: "session-123",
					}).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Nil request",
			req:  nil,
			mockCalls: func(t *testing.T, mockRepo *mocks.MockUserRepository) {
				// No mock calls expected
			},
			wantErr: response.ErrorUnAuthorized,
		},
		{
			name: "Empty ID in request",
			req: &models.AuthClaims{
				ID:      "",
				Subject: "user-456",
			},
			mockCalls: func(t *testing.T, mockRepo *mocks.MockUserRepository) {
				// No mock calls expected
			},
			wantErr: response.ErrorUnAuthorized,
		},
		{
			name: "Empty Subject in request",
			req: &models.AuthClaims{
				ID:      "session-123",
				Subject: "",
			},
			mockCalls: func(t *testing.T, mockRepo *mocks.MockUserRepository) {
				// No mock calls expected
			},
			wantErr: response.ErrorUnAuthorized,
		},
		{
			name: "Both ID and Subject empty",
			req: &models.AuthClaims{
				ID:      "",
				Subject: "",
			},
			mockCalls: func(t *testing.T, mockRepo *mocks.MockUserRepository) {
				// No mock calls expected
			},
			wantErr: response.ErrorUnAuthorized,
		},
		{
			name: "First delete session call fails",
			req: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			mockCalls: func(t *testing.T, mockRepo *mocks.MockUserRepository) {
				// First call fails
				mockRepo.EXPECT().
					DeleteSessionByID(gomock.Any(), models.DeleteSessionByIDParams{
						ID:     "session-123",
						UserID: "user-456",
					}).
					Return(errors.New("database error"))
				// No second call expected since first fails
			},
			wantErr: response.ErrorInternalServerError,
		},
		{
			name: "Second delete session call fails",
			req: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			mockCalls: func(t *testing.T, mockRepo *mocks.MockUserRepository) {
				// First call succeeds
				mockRepo.EXPECT().
					DeleteSessionByID(gomock.Any(), models.DeleteSessionByIDParams{
						ID:     "session-123",
						UserID: "user-456",
					}).
					Return(nil)

				// Second call fails
				mockRepo.EXPECT().
					DeleteSessionByID(gomock.Any(), models.DeleteSessionByIDParams{
						ID:     "user-456",
						UserID: "session-123",
					}).
					Return(errors.New("cache error"))
			},
			wantErr: response.ErrorInternalServerError,
		},
		{
			name: "Session not found (sql.ErrNoRows) in first call",
			req: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			mockCalls: func(t *testing.T, mockRepo *mocks.MockUserRepository) {
				// First call returns no rows found (handled in deleteSessionByID)
				mockRepo.EXPECT().
					DeleteSessionByID(gomock.Any(), models.DeleteSessionByIDParams{
						ID:     "session-123",
						UserID: "user-456",
					}).
					Return(sql.ErrNoRows)
				// No second call expected since first fails
			},
			wantErr: response.ErrorInternalServerError, // Logout converts any deleteSessionByID error to InternalServerError
		},
		{
			name: "Session not found (sql.ErrNoRows) in second call",
			req: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			mockCalls: func(t *testing.T, mockRepo *mocks.MockUserRepository) {
				// First call succeeds
				mockRepo.EXPECT().
					DeleteSessionByID(gomock.Any(), models.DeleteSessionByIDParams{
						ID:     "session-123",
						UserID: "user-456",
					}).
					Return(nil)

				// Second call returns no rows found (handled in deleteSessionByID)
				mockRepo.EXPECT().
					DeleteSessionByID(gomock.Any(), models.DeleteSessionByIDParams{
						ID:     "user-456",
						UserID: "session-123",
					}).
					Return(sql.ErrNoRows)
			},
			wantErr: response.ErrorInternalServerError, // Logout converts any deleteSessionByID error to InternalServerError
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockUserRepository(ctrl)
			tt.mockCalls(t, mockRepo)

			service := &UserService{
				repository: mockRepo,
				logger:     zap.NewNop(),
			}

			ctx := context.Background()
			err := service.Logout(ctx, tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test concurrent logout operations
func TestUserService_Logout_Concurrent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	// Set up expectations for multiple concurrent logout operations
	for i := 0; i < 10; i++ {
		mockRepo.EXPECT().
			DeleteSessionByID(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(2) // Each logout calls DeleteSessionByID twice
	}

	service := &UserService{
		repository: mockRepo,
		logger:     zap.NewNop(),
	}

	ctx := context.Background()
	concurrency := 10
	errChan := make(chan error, concurrency)

	// Launch concurrent logout operations
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			req := &models.AuthClaims{
				ID:      "session-" + string(rune(id)),
				Subject: "user-" + string(rune(id)),
			}
			err := service.Logout(ctx, req)
			errChan <- err
		}(i)
	}

	// Collect results
	for i := 0; i < concurrency; i++ {
		err := <-errChan
		assert.NoError(t, err)
	}
}

// Benchmark test for Logout function
func BenchmarkUserService_Logout(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	// Set up expectations for benchmark
	mockRepo.EXPECT().
		DeleteSessionByID(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	service := &UserService{
		repository: mockRepo,
		logger:     zap.NewNop(),
	}

	req := &models.AuthClaims{
		ID:      "benchmark-session",
		Subject: "benchmark-user",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.Logout(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Unit tests for generateToken function
func TestUserService_generateToken(t *testing.T) {
	// Valid test configuration
	validConfig := &configs.Config{
		App: configs.AppConfig{
			Name: "test-app",
		},
		Function: configs.FunctionConfig{
			Auth: configs.FunctionAuth{
				SecretKey: configs.FunctionAuthSecretKey{
					SessionID: "0123456789abcdef0123456789abcdef", // 32 bytes for session encryption
				},
				Token: configs.FunctionAuthToken{
					SecretKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", // Valid 64 hex chars
				},
			},
		},
	}

	validUser := models.User{
		ID:    "user-12345",
		Email: "test@example.com",
	}

	type testCase struct {
		name           string
		secretKey      string
		expiration     time.Duration
		user           models.User
		config         *configs.Config
		validateResult func(*testing.T, string, string, time.Time, error)
	}

	tests := []testCase{
		{
			name:       "Successful token generation",
			secretKey:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expiration: time.Hour,
			user:       validUser,
			config:     validConfig,
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)
				assert.NotEmpty(t, idToken)
				assert.False(t, expiredAt.IsZero())

				// Verify the token is valid PASETO format
				assert.Contains(t, tokenString, "v4.local.")

				// Verify expiration is approximately correct (within 1 second tolerance)
				expectedExpiration := time.Now().Add(time.Hour)
				assert.WithinDuration(t, expectedExpiration, expiredAt, time.Second)
			},
		},
		{
			name:       "Invalid secret key - too short",
			secretKey:  "short",
			expiration: time.Hour,
			user:       validUser,
			config:     validConfig,
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				assert.Empty(t, idToken)
				assert.True(t, expiredAt.IsZero())
			},
		},
		{
			name:       "Invalid secret key - non-hex characters",
			secretKey:  "gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg", // 64 chars but not hex
			expiration: time.Hour,
			user:       validUser,
			config:     validConfig,
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				assert.Empty(t, idToken)
				assert.True(t, expiredAt.IsZero())
			},
		},
		{
			name:       "Empty secret key",
			secretKey:  "",
			expiration: time.Hour,
			user:       validUser,
			config:     validConfig,
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				assert.Empty(t, idToken)
				assert.True(t, expiredAt.IsZero())
			},
		},
		{
			name:       "Zero expiration duration",
			secretKey:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expiration: 0,
			user:       validUser,
			config:     validConfig, // Function currently returns no error for zero expiration
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				// Current behavior: function returns empty values but no error for zero expiration
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				assert.Empty(t, idToken)
				assert.True(t, expiredAt.IsZero())
			},
		},
		{
			name:       "Negative expiration duration",
			secretKey:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expiration: -time.Hour,
			user:       validUser,
			config:     validConfig, // Function currently returns no error for negative expiration
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				// Current behavior: function returns empty values but no error for negative expiration
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				assert.Empty(t, idToken)
				assert.True(t, expiredAt.IsZero()) // Zero time when function returns early
			},
		},
		{
			name:       "Long expiration duration",
			secretKey:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expiration: 365 * 24 * time.Hour, // 1 year
			user:       validUser,
			config:     validConfig,
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)
				assert.NotEmpty(t, idToken)
				assert.False(t, expiredAt.IsZero())

				// Verify expiration is approximately 1 year from now
				expectedExpiration := time.Now().Add(365 * 24 * time.Hour)
				assert.WithinDuration(t, expectedExpiration, expiredAt, time.Second)
			},
		},
		{
			name:       "Invalid session ID encryption key",
			secretKey:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expiration: time.Hour,
			user:       validUser,
			config: &configs.Config{
				App: configs.AppConfig{
					Name: "test-app",
				},
				Function: configs.FunctionConfig{
					Auth: configs.FunctionAuth{
						SecretKey: configs.FunctionAuthSecretKey{
							SessionID: "invalid", // Too short for encryption
						},
					},
				},
			},
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				assert.Empty(t, idToken)
				assert.True(t, expiredAt.IsZero())
			},
		},
		{
			name:       "Empty app name in config",
			secretKey:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expiration: time.Hour,
			user:       validUser,
			config: &configs.Config{
				App: configs.AppConfig{
					Name: "", // Empty app name
				},
				Function: configs.FunctionConfig{
					Auth: configs.FunctionAuth{
						SecretKey: configs.FunctionAuthSecretKey{
							SessionID: "0123456789abcdef0123456789abcdef",
						},
					},
				},
			}, // Should still work, just with empty issuer
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)
				assert.NotEmpty(t, idToken)
				assert.False(t, expiredAt.IsZero())
			},
		},
		{
			name:       "Empty user ID",
			secretKey:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expiration: time.Hour,
			user: models.User{
				ID:    "", // Empty user ID
				Email: "test@example.com",
			},
			config: validConfig, // Should still work, just with empty subject
			validateResult: func(t *testing.T, tokenString, idToken string, expiredAt time.Time, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)
				assert.NotEmpty(t, idToken)
				assert.False(t, expiredAt.IsZero())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				config: tt.config,
				logger: zap.NewNop(),
			}

			tokenString, idToken, expiredAt, err := service.generateToken(tt.secretKey, tt.expiration, tt.user)

			tt.validateResult(t, tokenString, idToken, expiredAt, err)
		})
	}
}

// Test concurrent token generation
func TestUserService_generateToken_Concurrent(t *testing.T) {
	config := &configs.Config{
		App: configs.AppConfig{
			Name: "concurrent-test-app",
		},
		Function: configs.FunctionConfig{
			Auth: configs.FunctionAuth{
				SecretKey: configs.FunctionAuthSecretKey{
					SessionID: "0123456789abcdef0123456789abcdef",
				},
			},
		},
	}

	service := &UserService{
		config: config,
		logger: zap.NewNop(),
	}

	secretKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	expiration := time.Hour

	// Run multiple goroutines to test concurrent token generation
	const numGoroutines = 10
	results := make(chan error, numGoroutines)
	tokens := make(chan string, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			testUser := models.User{
				ID:    "user-" + string(rune(id)),
				Email: "user" + string(rune(id)) + "@example.com",
			}

			tokenString, idToken, expiredAt, err := service.generateToken(secretKey, expiration, testUser)
			if err != nil {
				results <- err
				tokens <- ""
				return
			}

			if tokenString == "" || idToken == "" || expiredAt.IsZero() {
				results <- assert.AnError
				tokens <- ""
				return
			}

			results <- nil
			tokens <- tokenString
		}(i)
	}

	// Check all results
	tokenSet := make(map[string]bool)
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		token := <-tokens
		assert.NoError(t, err)

		if token != "" {
			// Verify each token is unique
			assert.False(t, tokenSet[token], "Token should be unique")
			tokenSet[token] = true
		}
	}

	// Verify we got the expected number of unique tokens
	assert.Equal(t, numGoroutines, len(tokenSet))
}

// Benchmark token generation
func BenchmarkUserService_generateToken(b *testing.B) {
	config := &configs.Config{
		App: configs.AppConfig{
			Name: "benchmark-app",
		},
		Function: configs.FunctionConfig{
			Auth: configs.FunctionAuth{
				SecretKey: configs.FunctionAuthSecretKey{
					SessionID: "0123456789abcdef0123456789abcdef",
				},
			},
		},
	}

	service := &UserService{
		config: config,
		logger: zap.NewNop(),
	}

	secretKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	expiration := time.Hour
	user := models.User{
		ID:    "benchmark-user",
		Email: "benchmark@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, err := service.generateToken(secretKey, expiration, user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

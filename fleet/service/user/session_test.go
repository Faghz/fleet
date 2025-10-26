package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/elzestia/fleet/pkg/mocks"
	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test user service with Redis for session tests
func createTestUserServiceWithRedis(t *testing.T) (*UserService, *mocks.MockUserRepository) {
	service, mockRepo, _ := createTestUserService(t)

	// Use empty Redis client like in user_test.go
	redisClient := &redis.Client{}
	service.redis = redisClient

	return service, mockRepo
}

// Table-driven test for insertSessions
func TestUserService_insertSessions(t *testing.T) {
	tests := []struct {
		name          string
		session       models.InsertSessionParams
		mockError     error
		expectedError error
	}{
		{
			name: "Success",
			session: models.InsertSessionParams{
				ID: "session-123",
				ExpiresAt: pgtype.Timestamptz{
					Time:  time.Now().Add(time.Hour),
					Valid: true,
				},
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "Database error",
			session: models.InsertSessionParams{
				ID: "session-123",
				ExpiresAt: pgtype.Timestamptz{
					Time:  time.Now().Add(time.Hour),
					Valid: true,
				},
			},
			mockError:     errors.New("failed to insert session"),
			expectedError: errors.New("failed to insert session"),
		},
		{
			name: "High load scenario",
			session: models.InsertSessionParams{
				ID: "session-high-load-123",
				ExpiresAt: pgtype.Timestamptz{
					Time:  time.Now().Add(time.Hour),
					Valid: true,
				},
			},
			mockError:     errors.New("database connection pool exhausted"),
			expectedError: errors.New("database connection pool exhausted"),
		},
		{
			name: "Large expiration time",
			session: models.InsertSessionParams{
				ID: "session-with-long-expiration",
				ExpiresAt: pgtype.Timestamptz{
					Time:  time.Now().Add(time.Hour * 24 * 365),
					Valid: true,
				}, // 1 year
			},
			mockError:     nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _ := createTestUserService(t)
			ctx := context.Background()

			mockRepo.EXPECT().InsertSession(ctx, tt.session).Return(tt.mockError)

			err := service.insertSessions(ctx, tt.session)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test GetSessionByID function
func TestUserService_GetSessionByID(t *testing.T) {
	now := time.Now()
	validSession := models.Session{
		ID:     "session-123",
		UserID: "user-456",
		ExpiresAt: pgtype.Timestamptz{
			Time:  now.Add(time.Hour),
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
		CreatedBy: "system",
	}

	tests := []struct {
		name       string
		authClaims *models.AuthClaims
		setupMocks func(*mocks.MockUserRepository)
		expect     func(*testing.T, models.Session, error)
	}{
		{
			name: "Success - Database Hit",
			authClaims: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			setupMocks: func(mockRepo *mocks.MockUserRepository) {
				// Cache miss first
				mockRepo.EXPECT().GetSessionCache(
					context.Background(),
					"user-456",
					"session-123",
				).Return(models.Session{}, redis.Nil)

				// Database returns session successfully
				mockRepo.EXPECT().GetSessionByEntityId(
					context.Background(),
					models.GetSessionByEntityIdParams{
						ID:     "session-123",
						UserID: "user-456",
					},
				).Return(validSession, nil)
			},
			expect: func(t *testing.T, session models.Session, err error) {
				assert.NoError(t, err)
				assert.Equal(t, validSession, session)
			},
		},
		{
			name: "Success - Cache Hit",
			authClaims: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			setupMocks: func(mockRepo *mocks.MockUserRepository) {
				// Cache hit
				mockRepo.EXPECT().GetSessionCache(
					context.Background(),
					"user-456",
					"session-123",
				).Return(validSession, nil)
			},
			expect: func(t *testing.T, session models.Session, err error) {
				assert.NoError(t, err)
				assert.Equal(t, validSession, session)
			},
		},
		{
			name:       "Error - Nil AuthClaims",
			authClaims: nil,
			setupMocks: func(mockRepo *mocks.MockUserRepository) {
				// No mocks needed as function should return early
			},
			expect: func(t *testing.T, session models.Session, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorUnAuthorized, err)
				assert.Equal(t, models.Session{}, session)
			},
		},
		{
			name: "Error - Empty ID in AuthClaims",
			authClaims: &models.AuthClaims{
				ID:      "",
				Subject: "user-456",
			},
			setupMocks: func(mockRepo *mocks.MockUserRepository) {
				// No mocks needed as function should return early
			},
			expect: func(t *testing.T, session models.Session, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorUnAuthorized, err)
				assert.Equal(t, models.Session{}, session)
			},
		},
		{
			name: "Error - Empty Subject in AuthClaims",
			authClaims: &models.AuthClaims{
				ID:      "session-123",
				Subject: "",
			},
			setupMocks: func(mockRepo *mocks.MockUserRepository) {
				// No mocks needed as function should return early
			},
			expect: func(t *testing.T, session models.Session, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorUnAuthorized, err)
				assert.Equal(t, models.Session{}, session)
			},
		},
		{
			name: "Error - Session Cache Not Found",
			authClaims: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			setupMocks: func(mockRepo *mocks.MockUserRepository) {
				// Cache miss first
				mockRepo.EXPECT().GetSessionCache(
					context.Background(),
					"user-456",
					"session-123",
				).Return(models.Session{}, redis.Nil)

				// Database returns sql.ErrNoRows
				mockRepo.EXPECT().GetSessionByEntityId(
					context.Background(),
					models.GetSessionByEntityIdParams{
						ID:     "session-123",
						UserID: "user-456",
					},
				).Return(models.Session{}, sql.ErrNoRows)
			},
			expect: func(t *testing.T, session models.Session, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorUnAuthorized, err)
				assert.Equal(t, models.Session{}, session)
			},
		},
		{
			name: "Error - Session Not Found in Database",
			authClaims: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			setupMocks: func(mockRepo *mocks.MockUserRepository) {
				// Cache miss first
				mockRepo.EXPECT().GetSessionCache(
					context.Background(),
					"user-456",
					"session-123",
				).Return(models.Session{}, redis.Nil)

				// Database returns sql.ErrNoRows
				mockRepo.EXPECT().GetSessionByEntityId(
					context.Background(),
					models.GetSessionByEntityIdParams{
						ID:     "session-123",
						UserID: "user-456",
					},
				).Return(models.Session{}, sql.ErrNoRows)
			},
			expect: func(t *testing.T, session models.Session, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorUnAuthorized, err)
				assert.Equal(t, models.Session{}, session)
			},
		},
		{
			name: "Error - Database Connection Error",
			authClaims: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			setupMocks: func(mockRepo *mocks.MockUserRepository) {
				// Cache miss first
				mockRepo.EXPECT().GetSessionCache(
					context.Background(),
					"user-456",
					"session-123",
				).Return(models.Session{}, redis.Nil)

				// Database returns connection error
				mockRepo.EXPECT().GetSessionByEntityId(
					context.Background(),
					models.GetSessionByEntityIdParams{
						ID:     "session-123",
						UserID: "user-456",
					},
				).Return(models.Session{}, errors.New("database connection failed"))
			},
			expect: func(t *testing.T, session models.Session, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorInternalServerError, err)
				assert.Equal(t, models.Session{}, session)
			},
		},
		{
			name: "Error - Database Connection Error",
			authClaims: &models.AuthClaims{
				ID:      "session-123",
				Subject: "user-456",
			},
			setupMocks: func(mockRepo *mocks.MockUserRepository) {
				// Cache miss first
				mockRepo.EXPECT().GetSessionCache(
					context.Background(),
					"user-456",
					"session-123",
				).Return(models.Session{}, redis.Nil)

				// Database returns connection error
				mockRepo.EXPECT().GetSessionByEntityId(
					context.Background(),
					models.GetSessionByEntityIdParams{
						ID:     "session-123",
						UserID: "user-456",
					},
				).Return(models.Session{}, errors.New("database connection failed"))
			},
			expect: func(t *testing.T, session models.Session, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorInternalServerError, err)
				assert.Equal(t, models.Session{}, session)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestUserServiceWithRedis(t)
			ctx := context.Background()

			// Setup mocks based on test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Call the function under test
			session, err := service.GetSessionByID(ctx, tt.authClaims)
			tt.expect(t, session, err)
		})
	}
}

// Test GetSessionByID with concurrent requests without cache (stress test)
func TestUserService_GetSessionByID_Concurrent_NoCache(t *testing.T) {
	service, mockRepo, _ := createTestUserService(t)
	ctx := context.Background()

	authClaims := &models.AuthClaims{
		ID:      "session-concurrent-nocache",
		Subject: "user-concurrent-nocache",
	}

	validSession := models.Session{
		ID:     "session-concurrent-nocache",
		UserID: "user-concurrent-nocache",
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedBy: "system",
	}

	// Always return cache miss
	mockRepo.EXPECT().GetSessionCache(
		ctx,
		"user-concurrent-nocache",
		"session-concurrent-nocache",
	).Return(models.Session{}, redis.Nil).Times(200)

	// Mock the database call to return the session
	mockRepo.EXPECT().GetSessionByEntityId(
		context.Background(),
		models.GetSessionByEntityIdParams{
			ID:     "session-concurrent-nocache",
			UserID: "user-concurrent-nocache",
		},
	).Return(validSession, nil).Times(200)

	// Run multiple goroutines to test concurrent access
	const numGoroutines = 200
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			session, err := service.GetSessionByID(ctx, authClaims)
			if err != nil {
				results <- err
				return
			}
			if session.ID != validSession.ID {
				results <- errors.New("unexpected session ID")
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

// Test GetSessionByID with concurrent requests with cache (stress test)
func TestUserService_GetSessionByID_Concurrent_WithCache(t *testing.T) {
	service, mockRepo, _ := createTestUserService(t)
	ctx := context.Background()

	authClaims := &models.AuthClaims{
		ID:      "session-concurrent-cache",
		Subject: "user-concurrent-cache",
	}

	validSession := models.Session{
		ID:     "session-concurrent-cache",
		UserID: "user-concurrent-cache",
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedBy: "system",
	}

	// Always return cache hit
	mockRepo.EXPECT().GetSessionCache(
		ctx,
		"user-concurrent-cache",
		"session-concurrent-cache",
	).Return(validSession, nil).Times(200)

	// No database calls expected when cache hits

	// Run multiple goroutines to test concurrent access
	const numGoroutines = 200
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			session, err := service.GetSessionByID(ctx, authClaims)
			if err != nil {
				results <- err
				return
			}
			if session.ID != validSession.ID {
				results <- errors.New("unexpected session ID")
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

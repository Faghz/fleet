package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/mocks"
	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/repository"
	"github.com/elzestia/fleet/pkg/transport/http/request"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Test service creation
func TestCreateService(t *testing.T) {
	config := &configs.Config{}
	logger := zap.NewNop()
	repo := &repository.Repository{}
	redisClient := &redis.Client{}

	service := CreateService(config, logger, repo, redisClient)

	assert.NotNil(t, service)
	assert.Equal(t, config, service.config)
	assert.Equal(t, logger, service.logger)
	assert.NotNil(t, service.repository)
}

// Helper function to create a test user service with mocks from mocks folder
func createTestUserService(t *testing.T) (*UserService, *mocks.MockUserRepository, *mocks.MockUserRepository) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockUserRepository(ctrl)

	// Minimal config for testing
	config := &configs.Config{}
	logger := zap.NewNop()

	service := &UserService{
		config:     config,
		logger:     logger,
		repository: mockRepo,
	}

	return service, mockRepo, mockRepo
}

// MockTx is a simple mock implementation of pgx.Tx for testing
type MockTx struct{}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) { return nil, nil }
func (m *MockTx) Commit(ctx context.Context) error          { return nil }
func (m *MockTx) Rollback(ctx context.Context) error        { return nil }
func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults { return nil }
func (m *MockTx) LargeObjects() pgx.LargeObjects                               { return pgx.LargeObjects{} }
func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}
func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return nil
}
func (m *MockTx) Conn() *pgx.Conn { return nil }

// Table-driven test for GetUserByParams
func TestUserService_GetUserByParams(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		mockReturn     models.User
		mockError      error
		expectedError  error
		expectedUserID string
	}{
		{
			name:  "Success",
			email: "test@example.com",
			mockReturn: models.User{
				ID:    "user-123",
				Email: "test@example.com",
			},
			mockError:      nil,
			expectedError:  nil,
			expectedUserID: "user-123",
		},
		{
			name:           "User not found",
			email:          "nonexistent@example.com",
			mockReturn:     models.User{},
			mockError:      sql.ErrNoRows,
			expectedError:  response.ErrorUserDatabaseUserNotFound,
			expectedUserID: "",
		},
		{
			name:           "Database error",
			email:          "test@example.com",
			mockReturn:     models.User{},
			mockError:      errors.New("database connection error"),
			expectedError:  response.ErrorInternalServerError,
			expectedUserID: "",
		},
		{
			name:           "Empty email",
			email:          "",
			mockReturn:     models.User{},
			mockError:      sql.ErrNoRows,
			expectedError:  response.ErrorUserDatabaseUserNotFound,
			expectedUserID: "",
		},
		{
			name:  "Very long email",
			email: "very.long.email.address.that.exceeds.normal.expectations.and.tests.boundary.conditions.for.email.validation.and.database.constraints.which.might.cause.issues.in.some.systems.if.not.properly.handled.by.the.application@example.com",
			mockReturn: models.User{
				ID:    "user-with-long-email",
				Email: "very.long.email.address.that.exceeds.normal.expectations.and.tests.boundary.conditions.for.email.validation.and.database.constraints.which.might.cause.issues.in.some.systems.if.not.properly.handled.by.the.application@example.com",
			},
			mockError:      nil,
			expectedError:  nil,
			expectedUserID: "user-with-long-email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _ := createTestUserService(t)
			ctx := context.Background()

			mockRepo.EXPECT().GetUserByEmail(ctx, tt.email).Return(tt.mockReturn, tt.mockError)

			user, err := service.GetUserByEmail(ctx, tt.email)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Empty(t, user.ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserID, user.ID)
			}
		})
	}
}

// Table-driven test for RegisterUser
func TestUserService_RegisterUser(t *testing.T) {
	defaultConfig := &configs.Config{
		Function: configs.FunctionConfig{
			User: configs.FunctionUser{
				SecretKey: configs.FunctionUserSecretKey{
					Email:           "test-email-secret-key-exactly32b", // Exactly 32 bytes for AES-256
					EmailSalt:       "test-salt",
					EmailSaltLength: 16,
				},
			},
			Auth: configs.FunctionAuth{
				SecretKey: configs.FunctionAuthSecretKey{
					PasswordSalt: "ohirO31p28iP",
				},
			},
		},
	}

	type testCase struct {
		name      string
		request   *request.RegisterUserRequest
		config    *configs.Config
		setupMock func(context.Context, *testing.T, *mocks.MockUserRepository, *mocks.MockUserRepository, testCase)
		expect    func(context.Context, *testing.T, testCase, error)
	}

	tests := []testCase{
		{
			name: "Success - user registered successfully",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// Mock GetUserByEmail to return no user found (user doesn't exist)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, sql.ErrNoRows)

				// Mock successful transaction
				mockTx := &MockTx{}
				mockRepo.EXPECT().BeginTx(ctx).Return(mockTx, nil)

				// Create a mock Queries object for transaction
				ctrl := gomock.NewController(t)
				mockTxQueries := mocks.NewMockQuerier(ctrl)

				// Mock WithTx to return the mock transaction queries
				mockRepo.EXPECT().WithTx(mockTx).Return(mockTxQueries)

				// Mock successful InsertUser and InsertAuth on transaction queries
				mockTxQueries.EXPECT().InsertUser(ctx, gomock.Any()).Return(nil)
				mockTxQueries.EXPECT().InsertAuth(ctx, gomock.Any()).Return(nil)

				// Mock successful Commit
				mockRepo.EXPECT().CommitTx(mockTx).Return(nil)
			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Failure - user registration with long password",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123tolong72characterswhichisnotallowedbytheencryptionalgorithm",
				ConfirmPassword: "password123tolong72characterswhichisnotallowedbytheencryptionalgorithm",
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// Mock GetUserByEmail to return no user found (user doesn't exist)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, sql.ErrNoRows)

				// Mock successful transaction
				mockTx := &MockTx{}
				mockRepo.EXPECT().BeginTx(ctx).Return(mockTx, nil)

				// Create a mock Queries object for transaction
				ctrl := gomock.NewController(t)
				mockTxQueries := mocks.NewMockQuerier(ctrl)

				// Mock WithTx to return the mock transaction queries
				mockRepo.EXPECT().WithTx(mockTx).Return(mockTxQueries)

				// Mock successful InsertUser and InsertAuth on transaction queries
				mockTxQueries.EXPECT().InsertUser(ctx, gomock.Any()).Return(nil)
				MockUserRepository.EXPECT().RollbackTx(mockTx).Return(nil)

			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.Error(t, response.ErrorInternalServerError)
			},
		},
		{
			name: "Failure - failed insert user due to database error",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// Mock GetUserByEmail to return no user found (user doesn't exist)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, sql.ErrNoRows)

				// Mock successful transaction
				mockTx := &MockTx{}
				mockRepo.EXPECT().BeginTx(ctx).Return(mockTx, nil)

				// Create a mock Queries object for transaction
				ctrl := gomock.NewController(t)
				mockTxQueries := mocks.NewMockQuerier(ctrl)

				// Mock WithTx to return the mock transaction queries
				mockRepo.EXPECT().WithTx(mockTx).Return(mockTxQueries)

				// Mock successful InsertUser and InsertAuth on transaction queries
				mockTxQueries.EXPECT().InsertUser(ctx, gomock.Any()).Return(errors.New("database error"))
				MockUserRepository.EXPECT().RollbackTx(mockTx).Return(nil)

			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.Error(t, response.ErrorInternalServerError)
			},
		},
		{
			name: "Failure - failed insert auth due to database error",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// Mock GetUserByEmail to return no user found (user doesn't exist)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, sql.ErrNoRows)

				// Mock successful transaction
				mockTx := &MockTx{}
				mockRepo.EXPECT().BeginTx(ctx).Return(mockTx, nil)

				// Create a mock Queries object for transaction
				ctrl := gomock.NewController(t)
				mockTxQueries := mocks.NewMockQuerier(ctrl)

				// Mock WithTx to return the mock transaction queries
				mockRepo.EXPECT().WithTx(mockTx).Return(mockTxQueries)

				// Mock successful InsertUser and InsertAuth on transaction queries
				mockTxQueries.EXPECT().InsertUser(ctx, gomock.Any()).Return(nil)
				mockTxQueries.EXPECT().InsertAuth(ctx, gomock.Any()).Return(errors.New("database error"))

				MockUserRepository.EXPECT().RollbackTx(mockTx).Return(nil)

			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.Error(t, response.ErrorInternalServerError)
			},
		},
		{
			name: "Failure - failed commit due to database error",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// Mock GetUserByEmail to return no user found (user doesn't exist)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, sql.ErrNoRows)

				// Mock successful transaction
				mockTx := &MockTx{}
				mockRepo.EXPECT().BeginTx(ctx).Return(mockTx, nil)

				// Create a mock Queries object for transaction
				ctrl := gomock.NewController(t)
				mockTxQueries := mocks.NewMockQuerier(ctrl)

				// Mock WithTx to return the mock transaction queries
				mockRepo.EXPECT().WithTx(mockTx).Return(mockTxQueries)

				// Mock successful InsertUser and InsertAuth on transaction queries
				mockTxQueries.EXPECT().InsertUser(ctx, gomock.Any()).Return(nil)
				mockTxQueries.EXPECT().InsertAuth(ctx, gomock.Any()).Return(nil)

				// Mock Commit to return an error
				MockUserRepository.EXPECT().CommitTx(mockTx).Return(errors.New("commit error"))

				MockUserRepository.EXPECT().RollbackTx(mockTx).Return(nil)
			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.Error(t, response.ErrorInternalServerError)
			},
		},
		{
			name: "Failure - failed rollback due to database error",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// Mock GetUserByEmail to return no user found (user doesn't exist)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, sql.ErrNoRows)

				// Mock successful transaction
				mockTx := &MockTx{}
				mockRepo.EXPECT().BeginTx(ctx).Return(mockTx, nil)

				// Create a mock Queries object for transaction
				ctrl := gomock.NewController(t)
				mockTxQueries := mocks.NewMockQuerier(ctrl)

				// Mock WithTx to return the mock transaction queries
				mockRepo.EXPECT().WithTx(mockTx).Return(mockTxQueries)

				// Mock successful InsertUser and InsertAuth on transaction queries
				mockTxQueries.EXPECT().InsertUser(ctx, gomock.Any()).Return(nil)
				mockTxQueries.EXPECT().InsertAuth(ctx, gomock.Any()).Return(nil)

				// Mock Commit to return an error
				MockUserRepository.EXPECT().CommitTx(mockTx).Return(errors.New("commit error"))

				MockUserRepository.EXPECT().RollbackTx(mockTx).Return(errors.New("rollback error"))
			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.Error(t, response.ErrorInternalServerError)
			},
		},
		{
			name: "User email already exists",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// Mock GetUserByEmail to return existing user
				existingUser := models.User{
					ID:    "existing-user-123",
					Email: "encrypted-email",
				}
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(existingUser, nil)
			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorUserDatabaseUserEmailAlreadyUsed, err)
			},
		},
		{
			name: "Database error when checking existing user",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// Mock GetUserByEmail to return database error
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, errors.New("database connection error"))
			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorInternalServerError, err)
			},
		},
		{
			name: "Transaction begin failure",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// Mock GetUserByEmail to return no user found (user doesn't exist)
				mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(models.User{}, sql.ErrNoRows)
				// Mock BeginTx to return error
				MockUserRepository.EXPECT().BeginTx(ctx).Return(nil, errors.New("failed to begin transaction"))
			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorInternalServerError, err)
			},
		},
		{
			name: "Invalid email encryption config",
			request: &request.RegisterUserRequest{
				Name:            "John Doe",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			config: &configs.Config{
				Function: configs.FunctionConfig{
					User: configs.FunctionUser{
						SecretKey: configs.FunctionUserSecretKey{
							Email:           "short-key", // Invalid length for AES encryption
							EmailSalt:       "test-salt",
							EmailSaltLength: 16,
						},
					},
					Auth: configs.FunctionAuth{
						SecretKey: configs.FunctionAuthSecretKey{
							PasswordSalt: "ohirO31p28iP",
						},
					},
				},
			},
			setupMock: func(ctx context.Context, t *testing.T, mockRepo *mocks.MockUserRepository, MockUserRepository *mocks.MockUserRepository, tc testCase) {
				// No mocks needed as encryption will fail before any repository calls
			},
			expect: func(ctx context.Context, t *testing.T, tc testCase, err error) {
				assert.Error(t, err)
				assert.Equal(t, response.ErrorInternalServerError, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			service, mockRepo, MockUserRepository := createTestUserService(t)

			// Use provided config or default
			if tt.config != nil {
				service.config = tt.config
			} else {
				service.config = defaultConfig
			}

			tt.setupMock(ctx, t, mockRepo, MockUserRepository, tt)
			err := service.RegisterUser(ctx, tt.request)
			tt.expect(ctx, t, tt, err)
		})
	}
}

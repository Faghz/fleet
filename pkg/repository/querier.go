package repository

import (
	"context"

	"github.com/elzestia/fleet/pkg/models"
)

//go:generate mockgen -destination=../mocks/querier-mock.go -package=mocks github.com/elzestia/fleet/pkg/repository Querier
type Querier interface {
	DeleteSessionByID(ctx context.Context, arg models.DeleteSessionByIDParams) error
	GetAuthByUserUserID(ctx context.Context, userID string) (models.Auth, error)
	GetSessionByEntityId(ctx context.Context, arg models.GetSessionByEntityIdParams) (models.Session, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	InsertAuth(ctx context.Context, arg models.InsertAuthParams) error
	InsertSession(ctx context.Context, arg models.InsertSessionParams) error
	InsertUser(ctx context.Context, arg models.InsertUserParams) error
	UpdateUser(ctx context.Context, arg models.UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)

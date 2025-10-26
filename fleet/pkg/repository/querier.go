package repository

import (
	"context"

	"github.com/elzestia/fleet/pkg/models"
)

//go:generate mockgen -destination=../mocks/querier-mock.go -package=mocks github.com/elzestia/fleet/pkg/repository Querier
type Querier interface {
	// Vehicle operations
	GetVehicleByVehicleID(ctx context.Context, vehicleID string) (models.Vehicle, error)

	// Vehicle location operations
	InsertVehicleLocation(ctx context.Context, arg models.InsertVehicleLocationParams) (models.VehicleLocation, error)
	InsertVehicleLocationHistory(ctx context.Context, arg models.InsertVehicleLocationHistoryParams) (models.VehicleLocationHistory, error)

	// Auth operations
	GetAuthByUserUserID(ctx context.Context, userID string) (models.Auth, error)
	InsertAuth(ctx context.Context, arg models.InsertAuthParams) error

	// Session operations
	DeleteSessionByID(ctx context.Context, arg models.DeleteSessionByIDParams) error
	GetSessionByEntityId(ctx context.Context, arg models.GetSessionByEntityIdParams) (models.Session, error)
	InsertSession(ctx context.Context, arg models.InsertSessionParams) error

	// User operations
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	InsertUser(ctx context.Context, arg models.InsertUserParams) error
	UpdateUser(ctx context.Context, arg models.UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)

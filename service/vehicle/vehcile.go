package vehicle

import (
	"context"
	"errors"

	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/transport/mqtt/request"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// processVehicleLocation handles the business logic for vehicle location updates
func (s *VehicleService) ProcessVehicleLocationSync(ctx context.Context, req *request.VehicleLocationRequest) {
	// Acquire distributed lock for vehicle location update
	mutex := s.redis.Mutex.NewMutex("vehicle_location_mutex:" + req.VehicleID)
	if err := mutex.LockContext(ctx); err != nil {
		s.logger.Error("[ProcessVehicleLocationSync] Failed to acquire lock for vehicle location update", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		return
	}

	vehicle, err := s.repo.GetVehicleLatestLocationByVehicleID(ctx, req.VehicleID)
	if errors.Is(err, pgx.ErrNoRows) {
		s.logger.Warn("[ProcessVehicleLocationSync] Vehicle not found for location update", zap.String("vehicle_id", req.VehicleID))
		return
	}

	if err != nil {
		s.logger.Error("[ProcessVehicleLocationSync] Failed to get vehicle by vehicle number", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		return
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		s.logger.Error("[ProcessVehicleLocationSync] Failed to begin transaction for vehicle location update", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		return
	}

	defer func() {
		// Release the distributed lock
		if ok, err := mutex.UnlockContext(ctx); !ok || err != nil {
			s.logger.Error("[ProcessVehicleLocationSync] Failed to release lock for vehicle location update", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		}

		// Rollback transaction if not committed
		if err != nil {
			err = s.repo.RollbackTx(tx)
			if err != nil {
				s.logger.Error("[ProcessVehicleLocationSync] Failed to rollback transaction for vehicle location update", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
			}
		}
	}()

	s.logger.Debug("[ProcessVehicleLocationSync] Current vehicle timestamp", zap.String("vehicle_id", req.VehicleID), zap.Int64("current_timestamp", vehicle.Timestamp), zap.Int64("new_timestamp", req.Timestamp))
	// Update vehicle location only if the new timestamp is more recent
	if vehicle.Timestamp < req.Timestamp {
		params := models.InsertVehicleLocationParams{}
		err = params.ParseFromLocationSyncRequest(req, vehicle.EntityID)
		if err != nil {
			s.logger.Error("[ProcessVehicleLocationSync] Failed to parse vehicle location request", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
			return
		}

		_, err = s.repo.WithTx(tx).InsertVehicleLocation(ctx, params)
		if err != nil {
			s.logger.Error("[ProcessVehicleLocationSync] Failed to insert vehicle location", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
			return
		}
	}

	historyParams := models.InsertVehicleLocationHistoryParams{}
	err = historyParams.ParseFromLocationSyncRequest(req, vehicle.EntityID)
	if err != nil {
		s.logger.Error("[ProcessVehicleLocationSync] Failed to parse vehicle location history request", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		return
	}

	_, err = s.repo.WithTx(tx).InsertVehicleLocationHistory(ctx, historyParams)
	if err != nil {
		s.logger.Error("[ProcessVehicleLocationSync] Failed to insert vehicle location history", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("[ProcessVehicleLocationSync] Failed to commit transaction for vehicle location update", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		return
	}

	s.logger.Debug("[ProcessVehicleLocationSync] Successfully processed vehicle location",
		zap.String("vehicle_id", req.VehicleID),
		zap.Float64("latitude", req.Latitude),
		zap.Float64("longitude", req.Longitude),
		zap.Int64("timestamp", req.Timestamp))

}

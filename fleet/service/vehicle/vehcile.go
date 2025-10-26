package vehicle

import (
	"context"
	"errors"
	"time"

	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/transport/http/request"
	"github.com/elzestia/fleet/pkg/transport/http/response"
	mqttRequest "github.com/elzestia/fleet/pkg/transport/mqtt/request"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// processVehicleLocation handles the business logic for vehicle location updates
func (s *VehicleService) ProcessVehicleLocationSync(ctx context.Context, req *mqttRequest.VehicleLocationRequest) {
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

	err = s.checkAndPublishNearestPOI(ctx, vehicle.VehicleID, req.Latitude, req.Longitude, req.Timestamp)
	if err != nil {
		s.logger.Error("[ProcessVehicleLocationSync] Failed to check and publish nearest POI", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		return
	}

}

func (s *VehicleService) GetVehicleLatestLocationByVehicleID(ctx context.Context, vehicleID string) (*response.VehicleLocation, error) {
	data, err := s.repo.GetVehicleLatestLocationByVehicleID(ctx, vehicleID)
	if errors.Is(err, pgx.ErrNoRows) {
		s.logger.Debug("[GetVehicleLatestLocationByVehicleID] Vehicle location not found", zap.String("vehicle_id", vehicleID))
		return nil, response.ErrorVehicleNotFound
	}

	if err != nil {
		s.logger.Error("[GetVehicleLatestLocationByVehicleID] Failed to get vehicle latest location by vehicle ID", zap.String("vehicle_id", vehicleID), zap.Error(err))
		return nil, err
	}

	latitude, _ := data.Latitude.Int.Float64()
	longitude, _ := data.Longitude.Int.Float64()

	vehicleResp := response.VehicleLocation{
		VehicleID: vehicleID,
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: data.Timestamp,
		UpdatedAt: data.CreatedAt.Time.Format(time.RFC3339),
	}

	if data.UpdatedAt.Valid {
		vehicleResp.UpdatedAt = data.UpdatedAt.Time.Format(time.RFC3339)
	}

	return &vehicleResp, nil
}

func (s *VehicleService) GetVehicleLocationHistory(ctx context.Context, req *request.GetVehicleLocationHistoryRequest) (history []*response.VehicleLocation, err error) {
	vehicle, err := s.repo.GetVehicleByVehicleID(ctx, req.VehicleID)
	if errors.Is(err, pgx.ErrNoRows) {
		s.logger.Debug("[GetVehicleLocationHistory] Vehicle not found", zap.String("vehicle_id", req.VehicleID))
		return history, response.ErrorVehicleNotFound
	}

	if err != nil {
		s.logger.Error("[GetVehicleLocationHistory] Failed to get vehicle by vehicle ID", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		return history, err
	}

	records, err := s.repo.GetVehicleLocationHistoryByVehicleIDAndTimeRange(ctx, models.GetVehicleLocationHistoryByVehicleIDParams{
		VehicleID: vehicle.EntityID,
		Start:     req.Start,
		End:       req.End,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		s.logger.Debug("[GetVehicleLocationHistory] Vehicle location history not found", zap.String("vehicle_id", req.VehicleID))
		return history, nil
	}

	if err != nil {
		s.logger.Error("[GetVehicleLocationHistory] Failed to get vehicle location history by vehicle ID and time range", zap.String("vehicle_id", req.VehicleID), zap.Error(err))
		return history, err
	}

	for _, record := range records {
		latitude, _ := record.Latitude.Float64Value()
		longitude, _ := record.Longitude.Float64Value()

		history = append(history, &response.VehicleLocation{
			VehicleID: req.VehicleID,
			Latitude:  latitude.Float64,
			Longitude: longitude.Float64,
			Timestamp: record.Timestamp,
			UpdatedAt: record.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	return history, nil
}

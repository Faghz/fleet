package repository

import (
	"context"

	"github.com/elzestia/fleet/pkg/models"
)

const getVehicleByVehicleID = `-- name: GetVehicleByVehicleID :one
SELECT 
    entity_id,
    vehicle_id,
    vehicle_type,
    brand,
    model,
    year,
    status,
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
FROM vehicle
WHERE vehicle_id = $1 AND deleted_at IS NULL
LIMIT 1
`

func (q *Queries) GetVehicleByVehicleID(ctx context.Context, VehicleID string) (models.Vehicle, error) {
	row := q.db.QueryRow(ctx, getVehicleByVehicleID, VehicleID)
	var i models.Vehicle
	err := row.Scan(
		&i.EntityID,
		&i.VehicleID,
		&i.VehicleType,
		&i.Brand,
		&i.Model,
		&i.Year,
		&i.Status,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.DeletedAt,
		&i.DeletedBy,
	)
	return i, err
}

const insertVehicleLocation = `-- name: InsertVehicleLocation :one
INSERT INTO vehicle_location (
    entity_id,
    vehicle_entity_id,
    latitude,
    longitude,
    timestamp
) VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT (vehicle_entity_id) DO UPDATE SET
    latitude = EXCLUDED.latitude,
    longitude = EXCLUDED.longitude,
    timestamp = EXCLUDED.timestamp,
    updated_at = CURRENT_TIMESTAMP
RETURNING entity_id, vehicle_entity_id, latitude, longitude, timestamp, created_at, updated_at
`

func (q *Queries) InsertVehicleLocation(ctx context.Context, arg models.InsertVehicleLocationParams) (models.VehicleLocation, error) {
	row := q.db.QueryRow(ctx, insertVehicleLocation,
		arg.EntityID,
		arg.VehicleEntityID,
		arg.Latitude,
		arg.Longitude,
		arg.Timestamp,
	)
	var i models.VehicleLocation
	err := row.Scan(
		&i.EntityID,
		&i.VehicleEntityID,
		&i.Latitude,
		&i.Longitude,
		&i.Timestamp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertVehicleLocationHistory = `-- name: InsertVehicleLocationHistory :one
INSERT INTO vehicle_location_history (
    entity_id,
    vehicle_entity_id,
    latitude,
    longitude,
    timestamp
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING entity_id, vehicle_entity_id, latitude, longitude, timestamp, created_at, updated_at
`

func (q *Queries) InsertVehicleLocationHistory(ctx context.Context, arg models.InsertVehicleLocationHistoryParams) (models.VehicleLocationHistory, error) {
	row := q.db.QueryRow(ctx, insertVehicleLocationHistory,
		arg.EntityID,
		arg.VehicleEntityID,
		arg.Latitude,
		arg.Longitude,
		arg.Timestamp,
	)
	var i models.VehicleLocationHistory
	err := row.Scan(
		&i.EntityID,
		&i.VehicleEntityID,
		&i.Latitude,
		&i.Longitude,
		&i.Timestamp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVehicleLatestLocationByVehicleID = `-- name: GetVehicleLatestLocationByVehicleID :one
SELECT 
    v.entity_id,
    v.vehicle_id,
    vl.latitude,
    vl.longitude,
    vl.timestamp,
    vl.created_at,
    vl.updated_at
FROM vehicle_location vl
JOIN vehicle v ON v.entity_id = vl.vehicle_entity_id
WHERE v.vehicle_id = $1 AND v.deleted_at IS NULL
LIMIT 1
`

func (q *Queries) GetVehicleLatestLocationByVehicleID(ctx context.Context, vehicleID string) (models.GetVehicleLatestLocationByVehicleIDRow, error) {
	row := q.db.QueryRow(ctx, getVehicleLatestLocationByVehicleID, vehicleID)
	var i models.GetVehicleLatestLocationByVehicleIDRow
	err := row.Scan(
		&i.EntityID,
		&i.VehicleID,
		&i.Latitude,
		&i.Longitude,
		&i.Timestamp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

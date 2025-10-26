-- name: InsertVehicleLocationHistory :one
INSERT INTO vehicle_location_history (
    entity_id,
    vehicle_entity_id,
    latitude,
    longitude,
    timestamp
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: InsertVehicleLocation :one
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
RETURNING *;

-- name: GetVehicleLatestLocationByVehicleID :one
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
LIMIT 1;

-- name: GetVehicleByVehicleNumber :one
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
LIMIT 1;
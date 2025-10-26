-- name: GetPointOfInterests :many
SELECT 
    id,
    name,
    latitude,
    longitude,
    description,
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
FROM point_of_interest
WHERE deleted_at IS NULL;
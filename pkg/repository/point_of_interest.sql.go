package repository

import (
	"context"

	"github.com/elzestia/fleet/pkg/models"
)

const getPointOfInterests = `-- name: GetPointOfInterests :many
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
WHERE deleted_at IS NULL
`

func (q *Queries) GetPointOfInterests(ctx context.Context) ([]models.PointOfInterest, error) {
	rows, err := q.db.Query(ctx, getPointOfInterests)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.PointOfInterest{}
	for rows.Next() {
		var i models.PointOfInterest
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Latitude,
			&i.Longitude,
			&i.Description,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
			&i.DeletedAt,
			&i.DeletedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

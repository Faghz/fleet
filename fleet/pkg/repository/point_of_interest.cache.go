package repository

import (
	"context"
	"strconv"

	"github.com/elzestia/fleet/pkg/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

func buildHashKeyForPointOfInterest(name string) string {
	return "point_of_interest:" + name
}

func buildGeoKeyForPointOfInterest() string {
	return "point_of_interest"
}

func (r *Repository) SetGeoSpatialPointOfInterests(ctx context.Context, points []models.PointOfInterest) error {
	pipeline := r.redisConn.Pipeline()

	for _, point := range points {
		latitude, _ := point.Latitude.Float64Value()
		longitude, _ := point.Longitude.Float64Value()
		err := pipeline.GeoAdd(
			ctx,
			buildGeoKeyForPointOfInterest(),
			&redis.GeoLocation{
				Name:      point.Name,
				Longitude: longitude.Float64,
				Latitude:  latitude.Float64,
			},
		).Err()
		if err != nil {
			pipeline.Discard()
			return err
		}

		hashKey := buildHashKeyForPointOfInterest(point.Name)
		err = pipeline.HSet(
			ctx,
			hashKey,
			map[string]interface{}{
				"id":          point.ID,
				"name":        point.Name,
				"description": point.Description.String,
				"latitude":    latitude.Float64,
				"longitude":   longitude.Float64,
				"created_at":  point.CreatedAt.Time,
				"created_by":  point.CreatedBy.String,
				"updated_at":  point.UpdatedAt.Time,
				"updated_by":  point.UpdatedBy.String,
				"deleted_at":  point.DeletedAt.Time,
				"deleted_by":  point.DeletedBy.String,
			},
		).Err()
		if err != nil {
			pipeline.Discard()
			return err
		}
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetNearestPointOfInterests(ctx context.Context, latitude, longitude float64, radius float64, count int) ([]models.PointOfInterest, error) {
	geoQuery := &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "km",
		WithCoord: true,
		WithDist:  true,
		Count:     count,
		Sort:      "ASC",
	}

	locations, err := r.redisConn.GeoRadius(ctx, "point_of_interest", longitude, latitude, geoQuery).Result()
	if err != nil {
		return nil, err
	}

	points := make([]models.PointOfInterest, 0, len(locations))
	for _, loc := range locations {
		hashKey := buildHashKeyForPointOfInterest(loc.Name)
		hashData, err := r.redisConn.HGetAll(ctx, hashKey).Result()
		if err != nil {
			return nil, err
		}

		id, _ := strconv.Atoi(hashData["id"])
		latitudeFloat, _ := strconv.ParseFloat(hashData["latitude"], 64)
		longitudeFloat, _ := strconv.ParseFloat(hashData["longitude"], 64)

		var latitude, longitude pgtype.Numeric
		latitude.Scan(latitudeFloat)
		longitude.Scan(longitudeFloat)

		var description, createdBy, updatedBy, deletedBy pgtype.Text
		description.Scan(hashData["description"])
		createdBy.Scan(hashData["created_by"])
		updatedBy.Scan(hashData["updated_by"])
		deletedBy.Scan(hashData["deleted_by"])

		var createdAt, updatedAt, deletedAt pgtype.Timestamptz
		createdAt.Scan(hashData["created_at"])
		updatedAt.Scan(hashData["updated_at"])
		deletedAt.Scan(hashData["deleted_at"])

		pointOfInterest := models.PointOfInterest{
			ID:          int64(id),
			Name:        hashData["name"],
			Latitude:    latitude,
			Longitude:   longitude,
			Description: description,
			CreatedBy:   createdBy,
			CreatedAt:   createdAt,
			UpdatedBy:   updatedBy,
			UpdatedAt:   updatedAt,
			DeletedBy:   deletedBy,
			DeletedAt:   deletedAt,
		}

		points = append(points, pointOfInterest)
	}

	return points, nil
}

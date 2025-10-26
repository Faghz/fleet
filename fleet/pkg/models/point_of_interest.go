package models

import "github.com/jackc/pgx/v5/pgtype"

type GetPointOfInterestsRow struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Description pgtype.Text        `json:"description"`
	Latitude    pgtype.Numeric     `json:"latitude"`
	Longitude   pgtype.Numeric     `json:"longitude"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	CreatedBy   pgtype.Text        `json:"created_by"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	UpdatedBy   pgtype.Text        `json:"updated_by"`
	DeletedAt   pgtype.Timestamptz `json:"deleted_at"`
	DeletedBy   pgtype.Text        `json:"deleted_by"`
}

type GeoLocation struct {
	Name      string  `redis:"name"`
	Longitude float64 `redis:"longitude"`
	Latitude  float64 `redis:"latitude"`
}

package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Auth struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	Password  string             `json:"password"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	CreatedBy string             `json:"created_by"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	UpdatedBy pgtype.Text        `json:"updated_by"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	DeletedBy pgtype.Text        `json:"deleted_by"`
}

type PointOfInterest struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Latitude    pgtype.Numeric     `json:"latitude"`
	Longitude   pgtype.Numeric     `json:"longitude"`
	Description pgtype.Text        `json:"description"`
	CreatedBy   pgtype.Text        `json:"created_by"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedBy   pgtype.Text        `json:"updated_by"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	DeletedBy   pgtype.Text        `json:"deleted_by"`
	DeletedAt   pgtype.Timestamptz `json:"deleted_at"`
}

type SchemaMigration struct {
	Version string `json:"version"`
}

type Session struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	CreatedBy string             `json:"created_by"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	UpdatedBy pgtype.Text        `json:"updated_by"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	DeletedBy pgtype.Text        `json:"deleted_by"`
}

type User struct {
	ID        string             `json:"id"`
	Email     string             `json:"email"`
	EmailHash string             `json:"email_hash"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	CreatedBy string             `json:"created_by"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	UpdatedBy pgtype.Text        `json:"updated_by"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	DeletedBy pgtype.Text        `json:"deleted_by"`
}

type Vehicle struct {
	EntityID    pgtype.UUID        `json:"entity_id"`
	VehicleID   string             `json:"vehicle_id"`
	VehicleType pgtype.Text        `json:"vehicle_type"`
	Brand       pgtype.Text        `json:"brand"`
	Model       pgtype.Text        `json:"model"`
	Year        pgtype.Int4        `json:"year"`
	Status      string             `json:"status"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	CreatedBy   string             `json:"created_by"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	UpdatedBy   pgtype.Text        `json:"updated_by"`
	DeletedAt   pgtype.Timestamptz `json:"deleted_at"`
	DeletedBy   pgtype.Text        `json:"deleted_by"`
}

type VehicleLocation struct {
	EntityID        pgtype.UUID        `json:"entity_id"`
	VehicleEntityID pgtype.UUID        `json:"vehicle_entity_id"`
	Latitude        pgtype.Numeric     `json:"latitude"`
	Longitude       pgtype.Numeric     `json:"longitude"`
	Timestamp       int64              `json:"timestamp"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
}

type VehicleLocationHistory struct {
	EntityID        pgtype.UUID        `json:"entity_id"`
	VehicleEntityID pgtype.UUID        `json:"vehicle_entity_id"`
	Latitude        pgtype.Numeric     `json:"latitude"`
	Longitude       pgtype.Numeric     `json:"longitude"`
	Timestamp       int64              `json:"timestamp"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
}

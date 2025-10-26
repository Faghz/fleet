package models

import (
	"fmt"

	"github.com/elzestia/fleet/pkg/transport/mqtt/request"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type InsertVehicleLocationParams struct {
	EntityID        pgtype.UUID    `json:"entity_id"`
	VehicleEntityID pgtype.UUID    `json:"vehicle_entity_id"`
	Latitude        pgtype.Numeric `json:"latitude"`
	Longitude       pgtype.Numeric `json:"longitude"`
	Timestamp       int64          `json:"timestamp"`
}

func (p *InsertVehicleLocationParams) ParseFromLocationSyncRequest(req *request.VehicleLocationRequest, vehicleEntityId pgtype.UUID) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	var latitude, longitude pgtype.Numeric
	if err := latitude.Scan(fmt.Sprintf("%f", req.Latitude)); err != nil {
		return err
	}
	if err := longitude.Scan(fmt.Sprintf("%f", req.Longitude)); err != nil {
		return err
	}

	p.EntityID.Scan(id.String())
	p.VehicleEntityID.ScanUUID(vehicleEntityId)
	p.Latitude = latitude
	p.Longitude = longitude
	p.Timestamp = req.Timestamp

	return nil
}

type InsertVehicleLocationHistoryParams struct {
	EntityID        pgtype.UUID    `json:"entity_id"`
	VehicleEntityID pgtype.UUID    `json:"vehicle_entity_id"`
	Latitude        pgtype.Numeric `json:"latitude"`
	Longitude       pgtype.Numeric `json:"longitude"`
	Timestamp       int64          `json:"timestamp"`
}

func (p *InsertVehicleLocationHistoryParams) ParseFromLocationSyncRequest(req *request.VehicleLocationRequest, vehicleEntityId pgtype.UUID) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	var latitude, longitude pgtype.Numeric
	if err := latitude.Scan(fmt.Sprintf("%f", req.Latitude)); err != nil {
		return err
	}
	if err := longitude.Scan(fmt.Sprintf("%f", req.Longitude)); err != nil {
		return err
	}

	p.EntityID.Scan(id.String())
	p.VehicleEntityID.ScanUUID(vehicleEntityId)
	p.Latitude = latitude
	p.Longitude = longitude
	p.Timestamp = req.Timestamp

	return nil
}

type GetVehicleLatestLocationByVehicleIDRow struct {
	EntityID  pgtype.UUID        `json:"entity_id"`
	VehicleID string             `json:"vehicle_id"`
	Latitude  pgtype.Numeric     `json:"latitude"`
	Longitude pgtype.Numeric     `json:"longitude"`
	Timestamp int64              `json:"timestamp"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type GetVehicleLocationHistoryByVehicleIDParams struct {
	VehicleID pgtype.UUID `json:"vehicle_id"`
	Start     int64       `json:"start"`
	End       int64       `json:"end"`
}

type GetVehicleLocationHistoryByVehicleIDRow struct {
	EntityID  pgtype.UUID        `json:"entity_id"`
	Latitude  pgtype.Numeric     `json:"latitude"`
	Longitude pgtype.Numeric     `json:"longitude"`
	Timestamp int64              `json:"timestamp"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type ReachedNearestPointOfInterestEvent struct {
	VehicleID string                                `json:"vehicle_id"`
	Event     string                                `json:"event"`
	Location  ReachedNearestPointOfInterestLocation `json:"location"`
	Timestamp int64                                 `json:"timestamp"`
}

type ReachedNearestPointOfInterestLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

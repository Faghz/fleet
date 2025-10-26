package models

import "github.com/jackc/pgx/v5/pgtype"

type DeleteSessionByIDParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type GetSessionByEntityIdParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type InsertSessionParams struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
	CreatedBy string             `json:"created_by"`
	UpdatedBy pgtype.Text        `json:"updated_by"`
}

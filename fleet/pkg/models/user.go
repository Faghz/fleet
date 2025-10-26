package models

import "github.com/jackc/pgx/v5/pgtype"

type InsertUserParams struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	EmailHash string `json:"email_hash"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
}

type UpdateUserParams struct {
	Name      string      `json:"name"`
	UpdatedBy pgtype.Text `json:"updated_by"`
	ID        string      `json:"id"`
}

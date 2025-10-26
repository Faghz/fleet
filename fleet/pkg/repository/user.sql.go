package repository

import (
	"context"

	"github.com/elzestia/fleet/pkg/models"
)

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
    id,
    email,
    email_hash,
    name,
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
FROM "user"
WHERE (
    email_hash = $1
    AND deleted_at IS NULL
)
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i models.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.EmailHash,
		&i.Name,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.DeletedAt,
		&i.DeletedBy,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT
    id,
    email,
    email_hash,
    name,
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
FROM "user"
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (models.User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i models.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.EmailHash,
		&i.Name,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.DeletedAt,
		&i.DeletedBy,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :exec
INSERT INTO "user" (
    id,
    email,
    email_hash,
    name,
    created_by
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
`

func (q *Queries) InsertUser(ctx context.Context, arg models.InsertUserParams) error {
	_, err := q.db.Exec(ctx, insertUser,
		arg.ID,
		arg.Email,
		arg.EmailHash,
		arg.Name,
		arg.CreatedBy,
	)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE "user" SET
    name = $1,
    updated_by = $2
WHERE id = $3 AND deleted_at IS NULL
`

func (q *Queries) UpdateUser(ctx context.Context, arg models.UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.Name, arg.UpdatedBy, arg.ID)
	return err
}

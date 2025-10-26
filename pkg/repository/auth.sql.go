package repository

import (
	"context"

	"github.com/elzestia/fleet/pkg/models"
)

const getAuthByUserUserID = `-- name: GetAuthByUserUserID :one
SELECT
    id,
    user_id,
    password,
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
FROM auth
WHERE user_id = $1
`

func (q *Queries) GetAuthByUserUserID(ctx context.Context, userID string) (models.Auth, error) {
	row := q.db.QueryRow(ctx, getAuthByUserUserID, userID)
	var i models.Auth
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Password,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.DeletedAt,
		&i.DeletedBy,
	)
	return i, err
}

const insertAuth = `-- name: InsertAuth :exec
INSERT INTO auth (
    id,
    user_id,
    password,
    created_by
) VALUES (
    $1,
    $2,
    $3,
    $4
)
`

func (q *Queries) InsertAuth(ctx context.Context, arg models.InsertAuthParams) error {
	_, err := q.db.Exec(ctx, insertAuth,
		arg.ID,
		arg.UserID,
		arg.Password,
		arg.CreatedBy,
	)
	return err
}

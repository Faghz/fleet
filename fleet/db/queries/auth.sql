-- name: InsertAuth :exec
INSERT INTO auth (
    id,
    user_id,
    password,
    created_by
) VALUES (
    sqlc.arg('id'),
    sqlc.arg('user_id'),
    sqlc.arg('password'),
    sqlc.arg('created_by')
);

-- name: GetAuthByUserUserID :one
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
WHERE user_id = sqlc.arg('user_id');

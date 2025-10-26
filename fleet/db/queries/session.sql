-- name: InsertSession :exec
INSERT INTO session (
    id,
    user_id,
    expires_at,
    created_by,
    updated_by
) VALUES (
    sqlc.arg('id'),
    sqlc.arg('user_id'),
    sqlc.arg('expires_at'),
    sqlc.arg('created_by'),
    sqlc.arg('updated_by')
);


-- name: GetSessionByEntityId :one
SELECT
    id,
    user_id,
    expires_at,
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
FROM session
WHERE
    id = sqlc.arg('id')
    AND user_id = sqlc.arg('user_id')
    AND expires_at > NOW()
    AND deleted_at IS NULL
LIMIT 1;

-- name: DeleteSessionByID :exec
DELETE FROM session
WHERE
    id = sqlc.arg('id')
    AND user_id = sqlc.arg('user_id');

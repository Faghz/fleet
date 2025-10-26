-- name: GetUserByEmail :one
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
    email_hash = sqlc.arg('email')
    AND deleted_at IS NULL
);

-- name: GetUserByID :one
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
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;


-- name: InsertUser :exec
INSERT INTO "user" (
    id,
    email,
    email_hash,
    name,
    created_by
) VALUES (
    sqlc.arg('id'),
    sqlc.arg('email'),
    sqlc.arg('email_hash'),
    sqlc.arg('name'),
    sqlc.arg('created_by')
);

-- name: UpdateUser :exec
UPDATE "user" SET
    name = sqlc.arg('name'),
    updated_by = sqlc.arg('updated_by')
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

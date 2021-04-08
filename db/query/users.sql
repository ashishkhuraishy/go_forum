-- name: CreateUser :one
INSERT INTO users (full_name) 
VALUES ($1)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users 
SET full_name = $2
WHERE id = $1 
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;
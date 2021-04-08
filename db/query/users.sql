-- name: CreateUser :one
INSERT INTO users (full_name) 
VALUES ($1)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;


-- name: UpdateUser :exec
UPDATE users 
SET full_name = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM authors 
WHERE id = $1;
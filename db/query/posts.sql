-- name: CreatePost :one
INSERT INTO posts (
    user_id,
    title,
    descr
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListPostsFromUser :many
SELECT * FROM posts
WHERE user_id = $1
LIMIT $2
OFFSET $3;

-- name: UpdatePost :one
UPDATE posts 
SET title = $2,
descr = $3
WHERE id = $1 
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts 
WHERE id = $1;
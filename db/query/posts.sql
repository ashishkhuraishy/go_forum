-- name: CreatePost :one
INSERT INTO posts (
    user_id,
    title,
    content
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: ListPostsFromUser :many
SELECT * FROM posts
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: UpdatePost :one
UPDATE posts 
SET title = $2,
content = $3
WHERE id = $1 
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts 
WHERE id = $1;
-- name: AddLike :one
INSERT INTO likes (
    user_id,
    post_id,
    liked
) VALUES ($1, $2, $3)
RETURNING *;

-- name: CountPostLikes :one
SELECT SUM(
    CASE 
    WHEN liked 
    THEN 1 ELSE -1 
    END
)
FROM likes
WHERE post_id = $1;

-- name: GetLikeInfo :one
SELECT * FROM likes
WHERE user_id = $1
AND post_id = $2;

-- name: ListLikesOfUser :many
SELECT * FROM likes
WHERE user_id = $1
AND liked = true
ORDER BY id;

-- name: UpdateLike :one
UPDATE likes 
SET liked = $3
WHERE post_id = $1 
AND user_id = $2
RETURNING *;

-- name: DeleteLike :exec
DELETE FROM likes 
WHERE post_id = $1 
AND user_id = $2;
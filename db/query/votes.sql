-- name: AddVote :one
INSERT INTO votes (
    user_id,
    post_id,
    voted
) VALUES ($1, $2, $3)
RETURNING *;

-- name: CountVotesOfPost :one
SELECT SUM(
    CASE 
    WHEN voted 
    THEN 1 ELSE -1 
    END
)
FROM votes
WHERE post_id = $1;

-- name: GetVoteInfo :one
SELECT * FROM votes
WHERE user_id = $1
AND post_id = $2;

-- name: ListVotesOfUser :many
SELECT * FROM votes
WHERE user_id = $1
AND voted = true
ORDER BY id;

-- name: UpdateVote :one
UPDATE votes 
SET voted = $3
WHERE post_id = $1 
AND user_id = $2
RETURNING *;

-- name: DeleteVote :exec
DELETE FROM votes 
WHERE post_id = $1 
AND user_id = $2;
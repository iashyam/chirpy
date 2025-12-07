-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id) VALUES(
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: ListChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at DESC;

-- name: GetChipByID :one
SELECT * FROM chirps
WHERE id=$1;

-- name: DeleteChirpByID :exec
DELETE FROM chirps WHERE id=$1;

-- name: GetAuthorChrips :many
SELECT * from chirps
WHERE user_id=$1
ORDER by created_at ASC;
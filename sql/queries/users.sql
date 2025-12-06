-- name: AddUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password) VALUES(
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: Reset :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE
email=$1;

-- name: UpdateUser :one
UPDATE users
SET email=$2, hashed_password=$3, updated_at=Now()
WHERE id=$1
RETURNING *;

-- name: UpgradeUser :exec
UPDATE users
SET is_red=true WHERE id=$1;
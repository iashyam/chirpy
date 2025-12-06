-- name: CreateRefreshToken :one
INSERT INTO refreshtokens (token, 
created_at,
updated_at,
user_id,
expires_at) VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRefTokByID :one
SELECT * FROM refreshtokens WHERE token=$1 AND revoked_at IS NULL;

-- name: RevokeRefTok :exec
UPDATE refreshtokens
SET updated_at=$2, revoked_at=$2
WHERE token=$1
AND revoked_at IS NULL
AND expires_at>NOW();
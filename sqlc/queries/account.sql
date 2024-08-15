-- name: CreateAccount :one
INSERT INTO accounts (id, nickname, email, balance, created_at, updated_at, active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE active = true
ORDER BY created_at
LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET nickname = $2, email = $3, balance = $4, updated_at = $5, active = $6
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
UPDATE accounts
SET active = false
WHERE id = $1;

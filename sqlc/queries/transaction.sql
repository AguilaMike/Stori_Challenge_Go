-- name: CreateTransaction :one
INSERT INTO transactions (id, account_id, amount, type, input_file_id, input_date, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: ListTransactionsByAccount :many
SELECT * FROM transactions
WHERE account_id = $1
ORDER BY created_at
LIMIT $2 OFFSET $3;

-- name: GetTransactionSummary :one
SELECT
    SUM(CASE WHEN type = 'credit' THEN amount ELSE -amount END) as total_balance,
    COUNT(*) as transaction_count,
    AVG(CASE WHEN type = 'credit' THEN amount ELSE NULL END) as average_credit,
    AVG(CASE WHEN type = 'debit' THEN amount ELSE NULL END) as average_debit
FROM transactions
WHERE account_id = $1;

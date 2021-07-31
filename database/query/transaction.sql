--name: CreateTransaction :one
INSERT INTO transactions (user_id, account_id, category_id, name, transaction_type, amount, notes,date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

--name: UpdateTransaction :one
UPDATE transactions SET account_id = $2,
category_id = $3,
name = $4,
transaction_type = $5,
amount = $6,
notes = $7,
date = $8
WHERE transaction_id = $1
RETURNING *;

--name: GetTransactionByID :one
SELECT * FROM transactions
WHERE transaction_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1;

--name: ListTransactionsByUserID :many
SELECT * FROM transactions
WHERE user_id = $1
    AND deleted_at = '0001-01-01 00:00:00Z'
    AND date > $4
    AND date < $5
ORDER BY transaction_id
LIMIT  $2
OFFSET $3;

--name: ListTransactionsByAccountID :many
SELECT * FROM transactions
WHERE account_id = $1
  AND deleted_at = '0001-01-01 00:00:00Z'
  AND date > $4
  AND date < $5
ORDER BY transaction_id
LIMIT  $2
OFFSET $3;

--name: ListTransactionsByCategoryID :many
SELECT * FROM transactions
WHERE category_id = $1
  AND deleted_at = '0001-01-01 00:00:00Z'
  AND date > $4
  AND date < $5
ORDER BY transaction_id
LIMIT  $2
OFFSET $3;


--name: DeleteTransaction :one
UPDATE transactions SET deleted_at = now()
WHERE transaction_id = $1
  AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at


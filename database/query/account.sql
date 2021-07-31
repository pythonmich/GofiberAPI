--name: CreateAccount :one
INSERT INTO accounts(user_id, account_name, account_type, balance, currency)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


--name: UpdateAccount :one
UPDATE accounts SET balance = $2
WHERE account_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING *;

--name: GetAccountByID :one
SELECT * FROM accounts
WHERE account_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1;

--name: ListAccounts :many
SELECT * FROM accounts
WHERE user_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
ORDER BY account_id
LIMIT $2
OFFSET $3

--name: DeleteAccount :one
UPDATE accounts SET deleted_at = now()
WHERE account_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at;


--name: CreateMerchant :one
INSERT INTO merchant(user_id, name)
VALUES($1, $2)
RETURNING *;
--name: UpdateMerchant :one
UPDATE merchant SET name = $2
WHERE merchant_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING *;

--name: GetMerchantByID :one
SELECT * FROM merchant
WHERE merchant_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1;
--name: ListMerchants :one
SELECT * FROM merchant
WHERE user_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
ORDER BY merchant_id
LIMIT $2
OFFSET $3;

--name: DeleteMerchant :one
UPDATE merchant SET deleted_at = now()
WHERE merchant_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at;
--name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING *;

--name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1
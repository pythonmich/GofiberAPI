--name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING *;

--name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1;


--name: ListUsers :many
SELECT * FROM users
WHERE deleted_at = '0001-01-01 00:00:00Z'
ORDER BY user_id
LIMIT $1
OFFSET $2;

--name: DeleteUser :one
UPDATE users SET deleted_at = now(),
email = concat(email, '-DELETED-', uuid_generate_v4())
WHERE user_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at;


--name: UpdatePassword :one
UPDATE users SET password_changed_at = now(),
password_hash = $2
WHERE user_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING password_changed_at;

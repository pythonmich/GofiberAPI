--name: GrantRole :exec
INSERT INTO user_roles (user_id, role)
VALUES ($1, $2)
RETURNING *;


--name: RevokeRole :exec
DELETE FROM user_roles
WHERE user_id = $1
AND role = $2;


--name: GetUserRoleByID :one
SELECT role, created_at FROM user_roles
WHERE user_id = $1
LIMIT 1;




--name: ListUsersByRole :many
SELECT * FROM user_roles
WHERE role = $1
LIMIT $2
OFFSET $3
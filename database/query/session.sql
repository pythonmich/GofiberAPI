--name: SaveRefreshToken :exec
INSERT INTO sessions (user_id, device_id, refresh_token, expires_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, device_id) DO
    UPDATE
        SET refresh_token = $3,
            expires_at = $4,
            created_at = now();


--name: GetSession :one
SELECT * FROM sessions
WHERE user_id = $1
  AND device_id = $2
  AND refresh_token = $3
  AND to_timestamp(expires_at) >  now()
LIMIT 1;


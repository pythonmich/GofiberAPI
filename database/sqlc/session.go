package database

import (
	model "FiberFinanceAPI/database/models"
	"context"
)

const saveRefreshToken = `--name: SaveRefreshToken :exec
INSERT INTO sessions (user_id, device_id, refresh_token, expires_at)
VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id, device_id) 
	DO 
    	UPDATE 
        	SET refresh_token = $3,
            	expires_at = $4,
				created_at = now();
`

type SaveRefreshTokenParams struct {
	UserID       model.UserID   `json:"user_id"`
	DeviceID     model.DeviceID `json:"device_id"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresAt    int64          `json:"expires_at"`
}

func (q *Queries) SaveRefreshToken(ctx context.Context, args SaveRefreshTokenParams) error {
	q.logs.WithField("func", "database/sqlc/session.go -> SaveRefreshToken()").Debug()
	_, err := q.db.ExecContext(ctx, saveRefreshToken, args.UserID, args.DeviceID, args.RefreshToken, args.ExpiresAt)
	if err != nil {
		q.logs.WithError(err).Warn(err)
		return err
	}
	return err
}

type GetSessionsParams struct {
	UserID       model.UserID   `json:"user_id"`
	DeviceID     model.DeviceID `json:"device_id"`
	RefreshToken string         `json:"refresh_token"`
}

const getSession = `--name: GetSession :one
SELECT user_id, device_id, refresh_token, expires_at, created_at FROM sessions
WHERE user_id = $1
AND device_id = $2
AND refresh_token = $3
AND to_timestamp(expires_at) > now()
LIMIT 1`

func (q *Queries) GetSession(ctx context.Context, args GetSessionsParams) (model.Session, error) {
	q.logs.WithField("func", "database/sqlc/session.go -> GetSession()").Debug()
	row := q.db.QueryRowContext(ctx, getSession, args.UserID, args.DeviceID, args.RefreshToken)
	var session model.Session
	err := row.Scan(
		&session.UserID,
		&session.DeviceID,
		&session.RefreshToken,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		q.logs.WithError(err).Warn(err)
		return model.Session{}, err
	}
	return session, err
}

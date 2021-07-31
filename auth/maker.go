package auth

import "time"

// Maker is an interface for managing tokens to create and verify the TokenAccess created
type Maker interface {
	// CreateAccessToken  create a JWT Access Token and signs its before being issued to the user
	CreateAccessToken(userID string, duration time.Duration) (string, error)
	// VerifyAccessToken verifies the issued TokenAccess and returns the AccessPayload
	VerifyAccessToken(token string) (*AccessPayload, error)
	// CreateRefreshToken refreshes our token when it expires
	CreateRefreshToken(userID string, refreshDuration time.Duration) (string, error)
	// AccessTokenExpiresAt returns the period in which an access token expires at for our TokenAccess
	AccessTokenExpiresAt(token string) (int64, error)
	// VerifyRefreshToken verifies the issued RefreshToken and returns the RefreshPayload
	VerifyRefreshToken(token string) (*RefreshPayload, error)
	// RefreshTokenExpiresAt returns the period in which an access token expires at for our TokenAccess
	RefreshTokenExpiresAt(token string) (int64, error)
}

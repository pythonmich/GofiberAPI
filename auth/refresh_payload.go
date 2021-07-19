package auth

import "time"

// RefreshPayload  is the payload for our Refresh Token
type RefreshPayload struct {
	// SUB subject
	SUB string `json:"sub"`
	// EXP Expires At
	EXP int64 `json:"exp"`
}

// NewRefreshPayload creates a new refresh AccessPayload for our user
func NewRefreshPayload(userID string, duration time.Duration) (PayloadInterface, error) {
	return &RefreshPayload{
		SUB: userID,
		EXP: time.Now().Add(duration).Unix(),
	}, nil
}

// Valid Checks if AccessPayload is valid
func (p RefreshPayload) Valid() error {
	if time.Now().After(time.Unix(p.EXP, 0)) {
		return ErrExpiredToken
	}
	return nil
}

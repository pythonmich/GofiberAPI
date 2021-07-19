package auth

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type PayloadInterface interface {
	Valid() error
}

// AccessPayload is the payload for our TokenAccess
type AccessPayload struct {
	// JTI jwt ID
	JTI uuid.UUID `json:"jti"`
	// SUB subject
	SUB string `json:"sub"`
	// IAT Issued At
	IAT int64 `json:"iat"`
	// EXP Expires At
	EXP int64 `json:"exp"`
}

// NewAccessPayload creates a new payload for our user
func NewAccessPayload(userID string, duration time.Duration) (PayloadInterface, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &AccessPayload{
		JTI: tokenID,
		SUB: userID,
		IAT: time.Now().Unix(),
		EXP: time.Now().Add(duration).Unix(),
	}, nil
}

// Valid Checks if AccessPayload is valid
func (p AccessPayload) Valid() error {
	if time.Now().After(time.Unix(p.EXP, 0)) {
		return ErrExpiredToken
	}
	return nil
}

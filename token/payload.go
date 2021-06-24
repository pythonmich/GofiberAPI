package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("JWTToken is invalid")
	ErrExpiredToken = errors.New("JWTToken has expired")
)
// Payload is the payload for our JWTToken
type Payload struct {
	ID uuid.UUID `json:"id"`
	UserID string `json:"user_id"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}



// NewPayload creates a new payload for our user
func NewPayload(userID string, duration time.Duration) (*Payload,error) {
	tokenID, err := uuid.NewRandom(); if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (p Payload) Valid() error {
	if time.Now().After(p.ExpiresAt){
		return ErrExpiredToken
	}
	return nil
}
package token

import "time"

// Maker is an interface for managing tokens to create and verify the JWTToken created
type Maker interface {
	// CreateToken create a JWTToken and signs its before being issued to the user
	CreateToken(userID string, duration time.Duration) (string, error)
	// VerifyToken verifies the issued JWTToken and returns the Payload
	VerifyToken(token string) (*Payload,error)
}
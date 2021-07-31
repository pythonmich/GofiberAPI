package models

import "time"

type DeviceID string

// Session represents our User's session
type Session struct {
	UserID       UserID    `json:"user_id"`
	DeviceID     DeviceID  `json:"device_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    int64     `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// SessionDeviceID contains our device id
type SessionDeviceID struct {
	DeviceID DeviceID `json:"device_id" validate:"required"`
}

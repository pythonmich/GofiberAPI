package auth

import (
	"FiberFinanceAPI/utils"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const minimumSymmetricKeySize = 32

type TokenAccess struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	// We will return if the access token has expired
	AccessTokenExpiresAt  int64 `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt int64 `json:"-"` // we will store the time in our database
}

// JWTTokenMaker will implement our Maker interface to access its create and verity TokenAccess methods
type JWTTokenMaker struct {
	secretKey, refreshKey string
	logs                  *utils.StandardLogger
}

// NewJWTTokenMaker creates a new JWTTokenMaker
func NewJWTTokenMaker(secretKey, refreshKey string, logs *utils.StandardLogger) (Maker, error) {
	if len(secretKey) < minimumSymmetricKeySize || len(refreshKey) < minimumSymmetricKeySize {
		return nil, fmt.Errorf("invalid key size required %d characters", minimumSymmetricKeySize)
	}
	return &JWTTokenMaker{
		secretKey:  secretKey,
		refreshKey: refreshKey,
		logs:       logs,
	}, nil
}

// CreateAccessToken create a JWT Access Token and signs its before being issued to the user
func (J JWTTokenMaker) CreateAccessToken(userID string, duration time.Duration) (string, error) {
	J.logs.WithField("func", "auth/jwt_token.go -> TokenAccess()").Debug("Creating Token")
	payload, err := NewAccessPayload(userID, duration)
	if err != nil {
		J.logs.WithError(err).Warn("cannot create new payload")
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	J.logs.Info("Token Successfully Created")
	return jwtToken.SignedString([]byte(J.secretKey))
}

// VerifyAccessToken verifies the issued AccessToken and returns the AccessPayload
func (J JWTTokenMaker) VerifyAccessToken(token string) (*AccessPayload, error) {
	J.logs.WithField("func", "auth/jwt_token.go -> VerifyAccessToken()").Debug("Verifying Token")
	var claims AccessPayload
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		J.logs.WithField("func", "auth/jwt_token.go -> VerifyAccessToken() -> Anonymous func()").Debug("Verifying Token")
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			J.logs.WithError(ErrInvalidToken).Warn(ErrInvalidToken.Error())
			return nil, ErrInvalidToken
		}
		return []byte(J.secretKey), nil
	})
	if err != nil {
		//vErr validation error
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrExpiredToken) {
			J.logs.WithError(vErr).Warn("Token has expired")
			return &claims, ErrExpiredToken
		}
		J.logs.WithError(vErr).Warn(vErr.Error())
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*AccessPayload)
	if !ok {
		J.logs.WithError(ErrInvalidToken).Warn(ErrInvalidToken.Error())
		return nil, ErrInvalidToken
	}
	J.logs.Info("Token Verified")
	return payload, nil
}

// CreateRefreshToken refreshes our token when it expires
func (J JWTTokenMaker) CreateRefreshToken(userID string, refreshDuration time.Duration) (string, error) {
	J.logs.WithField("func", "auth/jwt_token.go -> CreateRefreshToken()").Debug("Refreshing Token")
	payload, err := NewRefreshPayload(userID, refreshDuration)
	if err != nil {
		J.logs.WithError(err).Warn("cannot create new refresh payload")
		return "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	J.logs.Debug("Refresh Token Successfully Created")
	return refreshToken.SignedString([]byte(J.refreshKey))
}

// AccessTokenExpiresAt returns the period in which an access token expires at for our TokenAccess
func (J JWTTokenMaker) AccessTokenExpiresAt(token string) (int64, error) {
	J.logs.WithField("func", "auth/jwt_token.go -> AccessTokenExpiresAt()").Debug("Access Token Expires At")
	payload, err := J.VerifyAccessToken(token)
	if err != nil {
		if errors.Is(err, ErrExpiredToken) {
			J.logs.WithError(err).Warn()
			return payload.EXP, err
		}
		J.logs.WithError(err).Warn(err.Error())
		return 0, err
	}
	return payload.EXP, nil
}

// VerifyRefreshToken verifies the issued RefreshToken and returns the RefreshPayload
func (J JWTTokenMaker) VerifyRefreshToken(token string) (*RefreshPayload, error) {
	J.logs.WithField("func", "auth/jwt_token.go -> VerifyRefreshToken()").Debug("Verifying Token")
	var claims RefreshPayload
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		J.logs.WithField("func", "auth/jwt_token.go -> VerifyRefreshToken() -> Anonymous func()").Debug("Verifying Token")
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			J.logs.WithError(ErrInvalidToken).Warn(ErrInvalidToken.Error())
			return nil, ErrInvalidToken
		}
		return []byte(J.refreshKey), nil
	})
	if err != nil {
		//vErr validation error
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrExpiredToken) {
			J.logs.WithError(vErr).Warn("Token has expired")
			return &claims, ErrExpiredToken
		}
		J.logs.WithError(vErr).Warn(vErr.Error())
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*RefreshPayload)
	if !ok {
		J.logs.WithError(ErrInvalidToken).Warn(ErrInvalidToken.Error())
		return nil, ErrInvalidToken
	}
	J.logs.Info("Token Verified")
	return payload, nil
}

// RefreshTokenExpiresAt returns the period in which a refresh token expires at
func (J JWTTokenMaker) RefreshTokenExpiresAt(token string) (int64, error) {
	J.logs.WithField("func", "auth/jwt_token.go -> RefreshTokenExpiresAt()").Debug("Refresh Token Expires At")
	payload, err := J.VerifyRefreshToken(token)
	if err != nil {
		if errors.Is(err, ErrExpiredToken) {
			J.logs.WithError(err).Warn(err.Error())
			return payload.EXP, err
		}
		J.logs.WithError(err).Warn(err.Error())
		return 0, err
	}
	return payload.EXP, nil
}

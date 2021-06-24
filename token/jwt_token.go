package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	logs "github.com/sirupsen/logrus"
	"time"
)


const minimumSymmetricKeySize = 32

// JWTTokenMaker will implement our Maker interface to access its create and verity JWTToken methods
type JWTTokenMaker struct {
	secretKey string
}



// NewJWTTokenMaker creates a new JWTTokenMaker
func NewJWTTokenMaker(secretKey string) (Maker,error) {
	if len(secretKey) < minimumSymmetricKeySize{
		return nil, fmt.Errorf("invalid key size required %d characters", minimumSymmetricKeySize)
	}
	return &JWTTokenMaker{
		secretKey: secretKey,
	}, nil
}
// CreateToken create a JWTToken and signs its before being issued to the user
func (J JWTTokenMaker) CreateToken(userID string, duration time.Duration) (string, error) {
	logs.WithField("func", "token/jwt_token.go -> CreateToken()").Debug("Creating Token")
	payload, err := NewPayload(userID, duration); if err != nil{
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(J.secretKey))
}
// VerifyToken verifies the issued JWTToken and returns the Payload
func (J JWTTokenMaker) VerifyToken(token string) (*Payload, error) {
	logs.WithField("func", "token/jwt_token.go -> VerifyToken()").Debug("Verifying Token")

	jwtToken, err := jwt.ParseWithClaims(token, Payload{}, func(token *jwt.Token) (interface{}, error) {
		logs.WithField("func", "token/jwt_token.go -> VerifyToken() -> Anonymous func()").Debug("Verifying Token")
		_, ok := token.Method.(*jwt.SigningMethodHMAC); if !ok{
			logs.WithError(ErrInvalidToken).Warn(ErrExpiredToken.Error())
			return nil, ErrInvalidToken
		}
		return []byte(J.secretKey),nil
	})
	if err != nil{
		//vErr validation error
		vErr, ok := err.(*jwt.ValidationError); if ok && errors.Is(vErr.Inner, ErrExpiredToken){
			logs.WithError(vErr).Warn(vErr.Error())
			return nil, ErrExpiredToken
		}
		logs.WithError(vErr).Warn(vErr.Error())
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload); if !ok{
		logs.WithError(ErrInvalidToken).Warn(ErrExpiredToken.Error())
		return nil, ErrInvalidToken
	}
	return payload, nil
}
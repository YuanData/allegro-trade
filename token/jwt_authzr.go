package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const keySize = 32

type JWTAuthzr struct {
	secretKey string
}

func NewJWTAuthzr(secretKey string) (Authzr, error) {
	if len(secretKey) < keySize {
		return nil, fmt.Errorf("key too short: minimum %d characters", keySize)
	}
	return &JWTAuthzr{secretKey}, nil
}

func (authzr *JWTAuthzr) CreateToken(membername string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(membername, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(authzr.secretKey))
	return token, payload, err
}

func (authzr *JWTAuthzr) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(authzr.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

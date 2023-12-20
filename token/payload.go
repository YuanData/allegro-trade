package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token detected")
	ErrExpiredToken = errors.New("expired token detected")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Membername  string    `json:"membername"`
	Role      string    `json:"role"`
	IssuedTime  time.Time `json:"issued_time"`
	ExpiredTime time.Time `json:"expired_time"`
}

func NewPayload(membername string, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Membername:  membername,
		Role:      role,
		IssuedTime:  time.Now(),
		ExpiredTime: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredTime) {
		return ErrExpiredToken
	}
	return nil
}

package token

import (
	"time"
)

type Authzr interface {
	CreateToken(membername string, role string, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}

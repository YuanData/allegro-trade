package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/YuanData/allegro-trade/util"
)

func TestJWTAuthzr(t *testing.T) {
	authzr, err := NewJWTAuthzr(util.RandomString(32))
	require.NoError(t, err)

	membername := util.RandomHolder()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := authzr.CreateToken(membername, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = authzr.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, membername, payload.Membername)
	require.WithinDuration(t, issuedAt, payload.IssuedTime, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredTime, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	authzr, err := NewJWTAuthzr(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := authzr.CreateToken(util.RandomHolder(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = authzr.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomHolder(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	authzr, err := NewJWTAuthzr(util.RandomString(32))
	require.NoError(t, err)

	payload, err = authzr.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

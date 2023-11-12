package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/YuanData/allegro-trade/util"
)

func createRandomMember(t *testing.T) Member {
	hashedPassword, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateMemberParams{
		Membername:       util.RandomHolder(),
		PasswordHash: hashedPassword,
		NameEntire:       util.RandomHolder(),
		Email:          util.RandomEmail(),
	}

	member, err := testQueries.CreateMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, member)

	require.Equal(t, arg.Membername, member.Membername)
	require.Equal(t, arg.PasswordHash, member.PasswordHash)
	require.Equal(t, arg.NameEntire, member.NameEntire)
	require.Equal(t, arg.Email, member.Email)
	require.True(t, member.PasswordChangedTime.IsZero())
	require.NotZero(t, member.CreatedTime)

	return member
}

func TestCreateMember(t *testing.T) {
	createRandomMember(t)
}

func TestGetMember(t *testing.T) {
	member1 := createRandomMember(t)
	member2, err := testQueries.GetMember(context.Background(), member1.Membername)
	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, member1.Membername, member2.Membername)
	require.Equal(t, member1.PasswordHash, member2.PasswordHash)
	require.Equal(t, member1.NameEntire, member2.NameEntire)
	require.Equal(t, member1.Email, member2.Email)
	require.WithinDuration(t, member1.PasswordChangedTime, member2.PasswordChangedTime, time.Second)
	require.WithinDuration(t, member1.CreatedTime, member2.CreatedTime, time.Second)
}

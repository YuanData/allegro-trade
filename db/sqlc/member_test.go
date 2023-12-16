package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

	member, err := testStore.CreateMember(context.Background(), arg)
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
	member2, err := testStore.GetMember(context.Background(), member1.Membername)
	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, member1.Membername, member2.Membername)
	require.Equal(t, member1.PasswordHash, member2.PasswordHash)
	require.Equal(t, member1.NameEntire, member2.NameEntire)
	require.Equal(t, member1.Email, member2.Email)
	require.WithinDuration(t, member1.PasswordChangedTime, member2.PasswordChangedTime, time.Second)
	require.WithinDuration(t, member1.CreatedTime, member2.CreatedTime, time.Second)
}

func TestUpdateMemberOnlyNameEntire(t *testing.T) {
	oldMember := createRandomMember(t)

	newNameEntire := util.RandomHolder()
	updatedMember, err := testStore.UpdateMember(context.Background(), UpdateMemberParams{
		Membername: oldMember.Membername,
		NameEntire: pgtype.Text{
			String: newNameEntire,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldMember.NameEntire, updatedMember.NameEntire)
	require.Equal(t, newNameEntire, updatedMember.NameEntire)
	require.Equal(t, oldMember.Email, updatedMember.Email)
	require.Equal(t, oldMember.PasswordHash, updatedMember.PasswordHash)
}

func TestUpdateMemberOnlyEmail(t *testing.T) {
	oldMember := createRandomMember(t)

	newEmail := util.RandomEmail()
	updatedMember, err := testStore.UpdateMember(context.Background(), UpdateMemberParams{
		Membername: oldMember.Membername,
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldMember.Email, updatedMember.Email)
	require.Equal(t, newEmail, updatedMember.Email)
	require.Equal(t, oldMember.NameEntire, updatedMember.NameEntire)
	require.Equal(t, oldMember.PasswordHash, updatedMember.PasswordHash)
}

func TestUpdateMemberOnlyPassword(t *testing.T) {
	oldMember := createRandomMember(t)

	newPassword := util.RandomString(8)
	newPasswordHash, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedMember, err := testStore.UpdateMember(context.Background(), UpdateMemberParams{
		Membername: oldMember.Membername,
		PasswordHash: pgtype.Text{
			String: newPasswordHash,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldMember.PasswordHash, updatedMember.PasswordHash)
	require.Equal(t, newPasswordHash, updatedMember.PasswordHash)
	require.Equal(t, oldMember.NameEntire, updatedMember.NameEntire)
	require.Equal(t, oldMember.Email, updatedMember.Email)
}

func TestUpdateMemberAllFields(t *testing.T) {
	oldMember := createRandomMember(t)

	newNameEntire := util.RandomHolder()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(8)
	newPasswordHash, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedMember, err := testStore.UpdateMember(context.Background(), UpdateMemberParams{
		Membername: oldMember.Membername,
		NameEntire: pgtype.Text{
			String: newNameEntire,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
		PasswordHash: pgtype.Text{
			String: newPasswordHash,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldMember.PasswordHash, updatedMember.PasswordHash)
	require.Equal(t, newPasswordHash, updatedMember.PasswordHash)
	require.NotEqual(t, oldMember.Email, updatedMember.Email)
	require.Equal(t, newEmail, updatedMember.Email)
	require.NotEqual(t, oldMember.NameEntire, updatedMember.NameEntire)
	require.Equal(t, newNameEntire, updatedMember.NameEntire)
}

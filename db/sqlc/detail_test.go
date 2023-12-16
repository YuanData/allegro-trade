package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/YuanData/allegro-trade/util"
)

func createRandomDetail(t *testing.T, trader Trader) Detail {
	arg := CreateDetailParams{
		TraderID: trader.ID,
		Number:    util.RandomAmount(),
	}

	detail, err := testStore.CreateDetail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, detail)

	require.Equal(t, arg.TraderID, detail.TraderID)
	require.Equal(t, arg.Number, detail.Number)

	require.NotZero(t, detail.ID)
	require.NotZero(t, detail.CreatedTime)

	return detail
}

func TestCreateDetail(t *testing.T) {
	trader := createRandomTrader(t)
	createRandomDetail(t, trader)
}

func TestGetDetail(t *testing.T) {
	trader := createRandomTrader(t)
	detail1 := createRandomDetail(t, trader)
	detail2, err := testStore.GetDetail(context.Background(), detail1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, detail2)

	require.Equal(t, detail1.ID, detail2.ID)
	require.Equal(t, detail1.TraderID, detail2.TraderID)
	require.Equal(t, detail1.Number, detail2.Number)
	require.WithinDuration(t, detail1.CreatedTime, detail2.CreatedTime, time.Second)
}

func TestListDetails(t *testing.T) {
	trader := createRandomTrader(t)
	for i := 0; i < 10; i++ {
		createRandomDetail(t, trader)
	}

	arg := ListDetailsParams{
		TraderID: trader.ID,
		Limit:     2,
		Offset:    2,
	}

	details, err := testStore.ListDetails(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, details, 2)

	for _, detail := range details {
		require.NotEmpty(t, detail)
		require.Equal(t, arg.TraderID, detail.TraderID)
	}
}

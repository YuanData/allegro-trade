package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/YuanData/allegro-trade/util"
)

func createRandomTrader(t *testing.T) Trader {
	member := createRandomMember(t)

	arg := CreateTraderParams{
		Holder:    member.Membername,
		Rest:  util.RandomAmount(),
		Symbol: util.RandomSymbol(),
	}

	trader, err := testQueries.CreateTrader(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trader)

	require.Equal(t, arg.Holder, trader.Holder)
	require.Equal(t, arg.Rest, trader.Rest)
	require.Equal(t, arg.Symbol, trader.Symbol)

	require.NotZero(t, trader.ID)
	require.NotZero(t, trader.CreatedTime)

	return trader
}

func TestCreateTrader(t *testing.T) {
	createRandomTrader(t)
}

func TestGetTrader(t *testing.T) {
	trader1 := createRandomTrader(t)
	trader2, err := testQueries.GetTrader(context.Background(), trader1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trader2)

	require.Equal(t, trader1.ID, trader2.ID)
	require.Equal(t, trader1.Holder, trader2.Holder)
	require.Equal(t, trader1.Rest, trader2.Rest)
	require.Equal(t, trader1.Symbol, trader2.Symbol)
	require.WithinDuration(t, trader1.CreatedTime, trader2.CreatedTime, time.Second)
}

func TestUpdateTrader(t *testing.T) {
	trader1 := createRandomTrader(t)

	arg := UpdateTraderParams{
		ID:      trader1.ID,
		Rest: util.RandomAmount(),
	}

	trader2, err := testQueries.UpdateTrader(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trader2)

	require.Equal(t, trader1.ID, trader2.ID)
	require.Equal(t, trader1.Holder, trader2.Holder)
	require.Equal(t, arg.Rest, trader2.Rest)
	require.Equal(t, trader1.Symbol, trader2.Symbol)
	require.WithinDuration(t, trader1.CreatedTime, trader2.CreatedTime, time.Second)
}

func TestDeleteTrader(t *testing.T) {
	trader1 := createRandomTrader(t)
	err := testQueries.DeleteTrader(context.Background(), trader1.ID)
	require.NoError(t, err)

	trader2, err := testQueries.GetTrader(context.Background(), trader1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, trader2)
}

func TestListTraders(t *testing.T) {
	var lastTrader Trader
	for i := 0; i < 4; i++ {
		lastTrader = createRandomTrader(t)
	}

	arg := ListTradersParams{
		Holder:  lastTrader.Holder,
		Limit:  2,
		Offset: 0,
	}

	traders, err := testQueries.ListTraders(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, traders)

	for _, trader := range traders {
		require.NotEmpty(t, trader)
		require.Equal(t, lastTrader.Holder, trader.Holder)
	}
}

package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/YuanData/allegro-trade/util"
)

func createRandomRecord(t *testing.T, trader1, trader2 Trader) Record {
	arg := CreateRecordParams{
		FromTraderID: trader1.ID,
		ToTraderID:   trader2.ID,
		Number:        util.RandomAmount(),
	}

	record, err := testStore.CreateRecord(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, record)

	require.Equal(t, arg.FromTraderID, record.FromTraderID)
	require.Equal(t, arg.ToTraderID, record.ToTraderID)
	require.Equal(t, arg.Number, record.Number)

	require.NotZero(t, record.ID)
	require.NotZero(t, record.CreatedTime)

	return record
}

func TestCreateRecord(t *testing.T) {
	trader1 := createRandomTrader(t)
	trader2 := createRandomTrader(t)
	createRandomRecord(t, trader1, trader2)
}

func TestGetRecord(t *testing.T) {
	trader1 := createRandomTrader(t)
	trader2 := createRandomTrader(t)
	record1 := createRandomRecord(t, trader1, trader2)

	record2, err := testStore.GetRecord(context.Background(), record1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, record2)

	require.Equal(t, record1.ID, record2.ID)
	require.Equal(t, record1.FromTraderID, record2.FromTraderID)
	require.Equal(t, record1.ToTraderID, record2.ToTraderID)
	require.Equal(t, record1.Number, record2.Number)
	require.WithinDuration(t, record1.CreatedTime, record2.CreatedTime, time.Second)
}

func TestListRecord(t *testing.T) {
	trader1 := createRandomTrader(t)
	trader2 := createRandomTrader(t)

	for i := 0; i < 5; i++ {
		createRandomRecord(t, trader1, trader2)
		createRandomRecord(t, trader2, trader1)
	}

	arg := ListRecordsParams{
		FromTraderID: trader1.ID,
		ToTraderID:   trader1.ID,
		Limit:         2,
		Offset:        2,
	}

	records, err := testStore.ListRecords(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, records, 2)

	for _, record := range records {
		require.NotEmpty(t, record)
		require.True(t, record.FromTraderID == trader1.ID || record.ToTraderID == trader1.ID)
	}
}

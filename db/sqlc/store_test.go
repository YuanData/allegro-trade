package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordTx(t *testing.T) {
	trader1 := createRandomTrader(t)
	trader2 := createRandomTrader(t)

	n := 6
	number := int64(800)

	errs := make(chan error)
	results := make(chan RecordTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.RecordTx(context.Background(), RecordTxParams{
				FromTraderID: trader1.ID,
				ToTraderID:   trader2.ID,
				Number:        number,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		record := result.Record
		require.NotEmpty(t, record)
		require.Equal(t, trader1.ID, record.FromTraderID)
		require.Equal(t, trader2.ID, record.ToTraderID)
		require.Equal(t, number, record.Number)
		require.NotZero(t, record.ID)
		require.NotZero(t, record.CreatedTime)

		_, err = testStore.GetRecord(context.Background(), record.ID)
		require.NoError(t, err)

		fromDetail := result.FromDetail
		require.NotEmpty(t, fromDetail)
		require.Equal(t, trader1.ID, fromDetail.TraderID)
		require.Equal(t, -number, fromDetail.Number)
		require.NotZero(t, fromDetail.ID)
		require.NotZero(t, fromDetail.CreatedTime)

		_, err = testStore.GetDetail(context.Background(), fromDetail.ID)
		require.NoError(t, err)

		toDetail := result.ToDetail
		require.NotEmpty(t, toDetail)
		require.Equal(t, trader2.ID, toDetail.TraderID)
		require.Equal(t, number, toDetail.Number)
		require.NotZero(t, toDetail.ID)
		require.NotZero(t, toDetail.CreatedTime)

		_, err = testStore.GetDetail(context.Background(), toDetail.ID)
		require.NoError(t, err)

		fromTrader := result.FromTrader
		require.NotEmpty(t, fromTrader)
		require.Equal(t, trader1.ID, fromTrader.ID)

		toTrader := result.ToTrader
		require.NotEmpty(t, toTrader)
		require.Equal(t, trader2.ID, toTrader.ID)


		diff1 := trader1.Rest - fromTrader.Rest
		diff2 := toTrader.Rest - trader2.Rest
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%number == 0)
		k := int(diff1 / number)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedTrader1, err := testStore.GetTrader(context.Background(), trader1.ID)
	require.NoError(t, err)

	updatedTrader2, err := testStore.GetTrader(context.Background(), trader2.ID)
	require.NoError(t, err)


	require.Equal(t, trader1.Rest-int64(n)*number, updatedTrader1.Rest)
	require.Equal(t, trader2.Rest+int64(n)*number, updatedTrader2.Rest)
}

func TestRecordTxDeadlock(t *testing.T) {
	trader1 := createRandomTrader(t)
	trader2 := createRandomTrader(t)

	n := 6
	number := int64(800)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromTraderID := trader1.ID
		toTraderID := trader2.ID

		if i%2 == 1 {
			fromTraderID = trader2.ID
			toTraderID = trader1.ID
		}

		go func() {
			_, err := testStore.RecordTx(context.Background(), RecordTxParams{
				FromTraderID: fromTraderID,
				ToTraderID:   toTraderID,
				Number:        number,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedTrader1, err := testStore.GetTrader(context.Background(), trader1.ID)
	require.NoError(t, err)

	updatedTrader2, err := testStore.GetTrader(context.Background(), trader2.ID)
	require.NoError(t, err)

	require.Equal(t, trader1.Rest, updatedTrader1.Rest)
	require.Equal(t, trader2.Rest, updatedTrader2.Rest)
}

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	RecordTx(ctx context.Context, arg RecordTxParams) (RecordTxResult, error)
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type RecordTxParams struct {
	FromTraderID int64 `json:"from_trader_id"`
	ToTraderID   int64 `json:"to_trader_id"`
	Number        int64 `json:"number"`
}

type RecordTxResult struct {
	Record    Record `json:"record"`
	FromTrader Trader  `json:"from_trader"`
	ToTrader   Trader  `json:"to_trader"`
	FromDetail   Detail    `json:"from_detail"`
	ToDetail     Detail    `json:"to_detail"`
}

func (store *SQLStore) RecordTx(ctx context.Context, arg RecordTxParams) (RecordTxResult, error) {
	var result RecordTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Record, err = q.CreateRecord(ctx, CreateRecordParams{
			FromTraderID: arg.FromTraderID,
			ToTraderID:   arg.ToTraderID,
			Number:        arg.Number,
		})
		if err != nil {
			return err
		}

		result.FromDetail, err = q.CreateDetail(ctx, CreateDetailParams{
			TraderID: arg.FromTraderID,
			Number:    -arg.Number,
		})
		if err != nil {
			return err
		}

		result.ToDetail, err = q.CreateDetail(ctx, CreateDetailParams{
			TraderID: arg.ToTraderID,
			Number:    arg.Number,
		})
		if err != nil {
			return err
		}

		if arg.FromTraderID < arg.ToTraderID {
			result.FromTrader, result.ToTrader, err = addAmount(ctx, q, arg.FromTraderID, -arg.Number, arg.ToTraderID, arg.Number)
		} else {
			result.ToTrader, result.FromTrader, err = addAmount(ctx, q, arg.ToTraderID, arg.Number, arg.FromTraderID, -arg.Number)
		}

		return err
	})

	return result, err
}

func addAmount(
	ctx context.Context,
	q *Queries,
	traderID1 int64,
	number1 int64,
	traderID2 int64,
	number2 int64,
) (trader1 Trader, trader2 Trader, err error) {
	trader1, err = q.AddTraderRest(ctx, AddTraderRestParams{
		ID:     traderID1,
		Number: number1,
	})
	if err != nil {
		return
	}

	trader2, err = q.AddTraderRest(ctx, AddTraderRestParams{
		ID:     traderID2,
		Number: number2,
	})
	return
}

package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/YuanData/allegro-trade/db/sqlc"
)

type recordRequest struct {
	FromTraderID int64  `json:"from_trader_id" binding:"required,min=1"`
	ToTraderID   int64  `json:"to_trader_id" binding:"required,min=1"`
	Number        int64  `json:"number" binding:"required,gt=0"`
	Symbol      string `json:"symbol" binding:"required,symbol"`
}

func (server *Server) createRecord(ctx *gin.Context) {
	var req recordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validTrader(ctx, req.FromTraderID, req.Symbol) {
		return
	}

	if !server.validTrader(ctx, req.ToTraderID, req.Symbol) {
		return
	}

	arg := db.RecordTxParams{
		FromTraderID: req.FromTraderID,
		ToTraderID:   req.ToTraderID,
		Number:        req.Number,
	}

	result, err := server.store.RecordTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validTrader(ctx *gin.Context, traderID int64, symbol string) bool {
	trader, err := server.store.GetTrader(ctx, traderID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if trader.Symbol != symbol {
		err := fmt.Errorf("trader [%d] symbol mismatch: %s vs %s", trader.ID, trader.Symbol, symbol)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}

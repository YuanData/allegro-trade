package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/token"
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

	fromTrader, valid := server.validTrader(ctx, req.FromTraderID, req.Symbol)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authztnPayloadKey).(*token.Payload)
	if fromTrader.Holder != authPayload.Membername {
		err := errors.New("from trader not under member")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validTrader(ctx, req.ToTraderID, req.Symbol)
	if !valid {
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

func (server *Server) validTrader(ctx *gin.Context, traderID int64, symbol string) (db.Trader, bool) {
	trader, err := server.store.GetTrader(ctx, traderID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return trader, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return trader, false
	}

	if trader.Symbol != symbol {
		err := fmt.Errorf("trader [%d] symbol mismatch: %s vs %s", trader.ID, trader.Symbol, symbol)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return trader, false
	}

	return trader, true
}

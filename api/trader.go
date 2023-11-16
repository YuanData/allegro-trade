package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/token"
)

type createTraderRequest struct {
	Symbol string `json:"symbol" binding:"required,symbol"`
}

func (server *Server) createTrader(ctx *gin.Context) {
	var req createTraderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authztnPayloadKey).(*token.Payload)
	arg := db.CreateTraderParams{
		Holder:    authPayload.Membername,
		Symbol: req.Symbol,
		Rest:  0,
	}

	trader, err := server.store.CreateTrader(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trader)
}

type getTraderRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTrader(ctx *gin.Context) {
	var req getTraderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	trader, err := server.store.GetTrader(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authztnPayloadKey).(*token.Payload)
	if trader.Holder != authPayload.Membername {
		err := errors.New("trader not under member")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trader)
}

type listTraderRequest struct {
	PageNum   int32 `form:"page_num" binding:"required,min=1"`
	PageLmt int32 `form:"page_lmt" binding:"required,min=5,max=10"`
}

func (server *Server) listTraders(ctx *gin.Context) {
	var req listTraderRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authztnPayloadKey).(*token.Payload)
	arg := db.ListTradersParams{
		Holder:  authPayload.Membername,
		Limit:  req.PageLmt,
		Offset: (req.PageNum - 1) * req.PageLmt,
	}

	traders, err := server.store.ListTraders(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, traders)
}

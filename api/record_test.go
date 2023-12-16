package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/YuanData/allegro-trade/db/mock"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/token"
	"github.com/YuanData/allegro-trade/util"
)

func TestRecordAPI(t *testing.T) {
	number := int64(10)

	member1, _ := randomMember(t)
	member2, _ := randomMember(t)
	member3, _ := randomMember(t)

	trader1 := randomTrader(member1.Membername)
	trader2 := randomTrader(member2.Membername)
	trader3 := randomTrader(member3.Membername)

	trader1.Symbol = util.ETH
	trader2.Symbol = util.ETH
	trader3.Symbol = util.BTC

	testCases := []struct {
		name          string
		body          gin.H
		setupAuthztn     func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"number":          number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(1).Return(trader2, nil)

				arg := db.RecordTxParams{
					FromTraderID: trader1.ID,
					ToTraderID:   trader2.ID,
					Number:        number,
				}
				store.EXPECT().RecordTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "MemberNotAuthorized",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"number":          number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member2.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(0)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthztnCredentials",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"number":          number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "FromTraderNotFound",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"number":          number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(db.Trader{}, db.ErrRecordNotFound)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(0)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToTraderNotFound",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"number":          number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(1).Return(db.Trader{}, db.ErrRecordNotFound)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromTraderSymbolMismatch",
			body: gin.H{
				"from_trader_id": trader3.ID,
				"to_trader_id":   trader2.ID,
				"number":          number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member3.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader3.ID)).Times(1).Return(trader3, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(0)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToTraderSymbolMismatch",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader3.ID,
				"number":          number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader3.ID)).Times(1).Return(trader3, nil)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidSymbol",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"number":          number,
				"symbol":        "NHK",
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeNumber",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"number":          -number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetTraderError",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"number":          number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Any()).Times(1).Return(db.Trader{}, sql.ErrConnDone)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "RecordTxError",
			body: gin.H{
				"from_trader_id": trader1.ID,
				"to_trader_id":   trader2.ID,
				"number":          number,
				"symbol":        util.ETH,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member1.Membername, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader1.ID)).Times(1).Return(trader1, nil)
				store.EXPECT().GetTrader(gomock.Any(), gomock.Eq(trader2.ID)).Times(1).Return(trader2, nil)
				store.EXPECT().RecordTx(gomock.Any(), gomock.Any()).Times(1).Return(db.RecordTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/records"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuthztn(t, request, server.tokenAuthzr)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

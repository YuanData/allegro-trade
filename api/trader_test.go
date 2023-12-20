package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

func TestGetTraderAPI(t *testing.T) {
	member, _ := randomMember(t)
	trader := randomTrader(member.Membername)

	testCases := []struct {
		name          string
		traderID     int64
		setupAuthztn     func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "Successful",
			traderID: trader.ID,
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Eq(trader.ID)).
					Times(1).
					Return(trader, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				verifyResponseBufferTrader(t, recorder.Body, trader)
			},
		},
		{
			name:      "MemberNotAuthorized",
			traderID: trader.ID,
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, "unauthorized_member", util.PrayerRole, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Eq(trader.ID)).
					Times(1).
					Return(trader, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthztnCredentials",
			traderID: trader.ID,
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Missing",
			traderID: trader.ID,
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Eq(trader.ID)).
					Times(1).
					Return(db.Trader{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "ServerFailure",
			traderID: trader.ID,
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Eq(trader.ID)).
					Times(1).
					Return(db.Trader{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "WrongID",
			traderID: 0,
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := fmt.Sprintf("/traders/%d", tc.traderID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuthztn(t, request, server.tokenAuthzr)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateTraderAPI(t *testing.T) {
	member, _ := randomMember(t)
	trader := randomTrader(member.Membername)

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
				"symbol": trader.Symbol,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTraderParams{
					Holder:    trader.Holder,
					Symbol: trader.Symbol,
					Rest:  0,
				}

				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(trader, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				verifyResponseBufferTrader(t, recorder.Body, trader)
			},
		},
		{
			name: "NoAuthztnCredentials",
			body: gin.H{
				"symbol": trader.Symbol,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ServerFailure",
			body: gin.H{
				"symbol": trader.Symbol,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Trader{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "WrongSymbol",
			body: gin.H{
				"symbol": "invalid",
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := "/traders"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuthztn(t, request, server.tokenAuthzr)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListTradersAPI(t *testing.T) {
	member, _ := randomMember(t)

	n := 8
	traders := make([]db.Trader, n)
	for i := 0; i < n; i++ {
		traders[i] = randomTrader(member.Membername)
	}

	type Query struct {
		PageNum   int
		PageLmt int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuthztn     func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful",
			query: Query{
				PageNum:   1,
				PageLmt: n,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTradersParams{
					Holder:  member.Membername,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(traders, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				verifyResponseBufferTraders(t, recorder.Body, traders)
			},
		},
		{
			name: "NoAuthztnCredentials",
			query: Query{
				PageNum:   1,
				PageLmt: n,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ServerFailure",
			query: Query{
				PageNum:   1,
				PageLmt: n,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Trader{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "WrongPageNum",
			query: Query{
				PageNum:   -1,
				PageLmt: n,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "WrongPageLmt",
			query: Query{
				PageNum:   1,
				PageLmt: 9999,
			},
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, member.Membername, member.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTraders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := "/traders"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_num", fmt.Sprintf("%d", tc.query.PageNum))
			q.Add("page_lmt", fmt.Sprintf("%d", tc.query.PageLmt))
			request.URL.RawQuery = q.Encode()

			tc.setupAuthztn(t, request, server.tokenAuthzr)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomTrader(holder string) db.Trader {
	return db.Trader{
		ID:       util.RandomInt(1, 2000),
		Holder:    holder,
		Rest:  util.RandomAmount(),
		Symbol: util.RandomSymbol(),
	}
}

func verifyResponseBufferTrader(t *testing.T, body *bytes.Buffer, trader db.Trader) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTrader db.Trader
	err = json.Unmarshal(data, &gotTrader)
	require.NoError(t, err)
	require.Equal(t, trader, gotTrader)
}

func verifyResponseBufferTraders(t *testing.T, body *bytes.Buffer, traders []db.Trader) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTraders []db.Trader
	err = json.Unmarshal(data, &gotTraders)
	require.NoError(t, err)
	require.Equal(t, traders, gotTraders)
}

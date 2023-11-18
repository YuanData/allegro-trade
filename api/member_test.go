package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/YuanData/allegro-trade/db/mock"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/util"
)

type eqCreateMemberParamsMatcher struct {
	arg      db.CreateMemberParams
	password string
}

func (e eqCreateMemberParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateMemberParams)
	if !ok {
		return false
	}

	err := util.VerifyPassword(e.password, arg.PasswordHash)
	if err != nil {
		return false
	}

	e.arg.PasswordHash = arg.PasswordHash
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateMemberParamsMatcher) String() string {
	return fmt.Sprintf("%v %v", e.arg, e.password)
}

func EqCreateMemberParams(arg db.CreateMemberParams, password string) gomock.Matcher {
	return eqCreateMemberParamsMatcher{arg, password}
}

func TestCreateMemberAPI(t *testing.T) {
	member, password := randomMember(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful",
			body: gin.H{
				"membername":  member.Membername,
				"password":  password,
				"name_entire": member.NameEntire,
				"email":     member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateMemberParams{
					Membername: member.Membername,
					NameEntire: member.NameEntire,
					Email:    member.Email,
				}
				store.EXPECT().
					CreateMember(gomock.Any(), EqCreateMemberParams(arg, password)).
					Times(1).
					Return(member, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMember(t, recorder.Body, member)
			},
		},
		{
			name: "ServerIssue",
			body: gin.H{
				"membername":  member.Membername,
				"password":  password,
				"name_entire": member.NameEntire,
				"email":     member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Member{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "MemberNameAlreadyUsed",
			body: gin.H{
				"membername":  member.Membername,
				"password":  password,
				"name_entire": member.NameEntire,
				"email":     member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Member{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "WrongMemberName",
			body: gin.H{
				"membername":  "wrong-member-name#",
				"password":  password,
				"name_entire": member.NameEntire,
				"email":     member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "WrongEmail",
			body: gin.H{
				"membername":  member.Membername,
				"password":  password,
				"name_entire": member.NameEntire,
				"email":     "wrong-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "PasswordUnderlength",
			body: gin.H{
				"membername":  member.Membername,
				"password":  "808",
				"name_entire": member.NameEntire,
				"email":     member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
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

			url := "/members"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginMemberAPI(t *testing.T) {
	member, password := randomMember(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful",
			body: gin.H{
				"membername": member.Membername,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMember(gomock.Any(), gomock.Eq(member.Membername)).
					Times(1).
					Return(member, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "MemberNotFound",
			body: gin.H{
				"membername": "NotFound",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMember(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Member{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"membername": member.Membername,
				"password": "incorrect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMember(gomock.Any(), gomock.Eq(member.Membername)).
					Times(1).
					Return(member, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ServerIssue",
			body: gin.H{
				"membername": member.Membername,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMember(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Member{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "WrongMemberName",
			body: gin.H{
				"membername":  "wrong-member-name#",
				"password":  password,
				"name_entire": member.NameEntire,
				"email":     member.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMember(gomock.Any(), gomock.Any()).
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

			url := "/members/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomMember(t *testing.T) (member db.Member, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	member = db.Member{
		Membername:       util.RandomHolder(),
		PasswordHash: hashedPassword,
		NameEntire:       util.RandomHolder(),
		Email:          util.RandomEmail(),
	}
	return
}

func requireBodyMatchMember(t *testing.T, body *bytes.Buffer, member db.Member) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotMember db.Member
	err = json.Unmarshal(data, &gotMember)

	require.NoError(t, err)
	require.Equal(t, member.Membername, gotMember.Membername)
	require.Equal(t, member.NameEntire, gotMember.NameEntire)
	require.Equal(t, member.Email, gotMember.Email)
	require.Empty(t, gotMember.PasswordHash)
}

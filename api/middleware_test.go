package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/YuanData/allegro-trade/token"
)

func addAuthztn(
	t *testing.T,
	request *http.Request,
	tokenAuthzr token.Authzr,
	authztnType string,
	membername string,
	duration time.Duration,
) {
	token, payload, err := tokenAuthzr.CreateToken(membername, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authztnType, token)
	request.Header.Set(authztnHeaderKey, authorizationHeader)
}

func TestAuthztnMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuthztn     func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful",
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, "member", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthztnCredentials",
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthztn",
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, "invalid", "member", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "IncorrectAuthztnFormat",
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, "", "member", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "TokenExpirationOccurred",
			setupAuthztn: func(t *testing.T, request *http.Request, tokenAuthzr token.Authzr) {
				addAuthztn(t, request, tokenAuthzr, authztnTypeBearer, "member", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authztnPath := "/authztn"
			server.router.GET(
				authztnPath,
				authztnMiddleware(server.tokenAuthzr),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authztnPath, nil)
			require.NoError(t, err)

			tc.setupAuthztn(t, request, server.tokenAuthzr)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

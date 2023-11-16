package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/YuanData/allegro-trade/token"
)

const (
	authztnHeaderKey  = "authorization"
	authztnTypeBearer = "bearer"
	authztnPayloadKey = "auth_payload"
)

func authztnMiddleware(tokenAuthzr token.Authzr) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authztnHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("header authoztn absent")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("wrong header authoztn format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authztnType := strings.ToLower(fields[0])
		if authztnType != authztnTypeBearer {
			err := fmt.Errorf("invalid authoztn type %s", authztnType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenAuthzr.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authztnPayloadKey, payload)
		ctx.Next()
	}
}

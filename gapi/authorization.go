package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/YuanData/allegro-trade/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeMember(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata err")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("no header err")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("wrong header authoztn format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("invalid authoztn type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.tokenAuthzr.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("access token err: %s", err)
	}

	if !hasPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}

	return payload, nil
}

func hasPermission(memberRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if memberRole == role {
			return true
		}
	}
	return false
}

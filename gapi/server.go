package gapi

import (
	"fmt"

	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/pb"
	"github.com/YuanData/allegro-trade/token"
	"github.com/YuanData/allegro-trade/util"
)

type Server struct {
	pb.UnimplementedAllegroTradeServer
	config     util.Config
	store      db.Store
	tokenAuthzr token.Authzr
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenAuthzr, err := token.NewJWTAuthzr(config.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("token authzr err: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenAuthzr: tokenAuthzr,
	}

	return server, nil
}

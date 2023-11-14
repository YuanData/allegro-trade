package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/token"
	"github.com/YuanData/allegro-trade/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenAuthzr token.Authzr
	router     *gin.Engine
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("symbol", validSymbol)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/members", server.createMember)
	router.POST("/members/login", server.loginMember)

	router.POST("/traders", server.createTrader)
	router.GET("/traders/:id", server.getTrader)
	router.GET("/traders", server.listTraders)

	router.POST("/records", server.createRecord)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"err resp": err.Error()}
}

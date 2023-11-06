package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/YuanData/allegro-trade/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/traders", server.createTrader)
	router.GET("/traders/:id", server.getTrader)
	router.GET("/traders", server.listTrader)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"err resp": err.Error()}
}

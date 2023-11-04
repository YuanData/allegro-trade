package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/YuanData/allegro-trade/api"
	db "github.com/YuanData/allegro-trade/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:pasd@localhost:5432/allegro_trade?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("sql open err: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("server start err", err)
	}
}

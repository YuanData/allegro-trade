package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/YuanData/allegro-trade/api"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("load config err: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("sql open err: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("server start err: ", err)
	}
}

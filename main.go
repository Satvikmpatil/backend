package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/satvikmpatil/simplebank/api"
	db "github.com/satvikmpatil/simplebank/db/sqlc"
	"github.com/satvikmpatil/simplebank/util"
)



func main() {
	config, err := util.LoadConfig(".")
	if err != nil{
		log.Fatal("cant load config ",err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("DB error", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil{
		log.Fatal("Canot start server:", err)
	}
}

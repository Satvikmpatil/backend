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

// ackage main

// import (
// 	"database/sql"
// 	"log"

// 	"github.com/Satvikmpatil/tolist/api"
// 	db "github.com/Satvikmpatil/tolist/db/sqlc"
// 	_ "github.com/lib/pq"
// )

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://postgres:mypassword@localhost:5500/postgres?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )

// func main() {
// 	conn, err := sql.Open(dbDriver, dbSource)
// 	if err != nil {
// 		log.Fatal("cannot connect to db:", err)
// 	}

// 	store := db.New(conn)
// 	server := api.NewServer(store)

// 	err = server.Start(serverAddress)
// 	if err != nil {
// 		log.Fatal("cannot start server:", err)
// 	}
// }
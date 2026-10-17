package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/satvikmpatil/simplebank/api"
	db "github.com/satvikmpatil/simplebank/db/sqlc"
	"github.com/satvikmpatil/simplebank/gapi"
	"github.com/satvikmpatil/simplebank/pb"
	"github.com/satvikmpatil/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)



func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cant load config ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("DB error", err)
	}
	store := db.NewStore(conn)

	// Run gRPC server (comment out runGinServer if you only want gRPC)
	runGrpcServer(config, store)

	// Or run HTTP server
	// runGinServer(config, store)
}

func runGrpcServer(config util.Config,store db.Store){
	server, err := gapi.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot crete server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil{
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil{
		log.Fatal("Cannot start gRPC server")
	}

}

func runGinServer(config util.Config,store db.Store){
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}
	err = server.Start(config.HHTPServerAddress)
	if err != nil{
		log.Fatal("Canot start server:", err)
	}
}
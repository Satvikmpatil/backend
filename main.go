package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/satvikmpatil/simplebank/api"
	db "github.com/satvikmpatil/simplebank/db/sqlc"
	_ "github.com/satvikmpatil/simplebank/doc/statik"
	"github.com/satvikmpatil/simplebank/gapi"
	"github.com/satvikmpatil/simplebank/pb"
	"github.com/satvikmpatil/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	// Run DB migrations
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	// Run gRPC server (comment out runGinServer if you only want gRPC)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)

	// Or run HTTP server
	// runGinServer(config, store)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}

	log.Println("db migrated successfully")
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

func runGatewayServer(config util.Config,store db.Store){
	server, err := gapi.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot crete server:", err)
	}

	
	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimplebankHandlerServer(ctx, grpcMux, server)

	if err != nil{
		log.Fatal("cannot create listener")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// Serve swagger documentation using statik
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs:", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HHTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Cannot start Gateway server")
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
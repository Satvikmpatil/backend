package gapi

import (
	db "github.com/satvikmpatil/simplebank/db/sqlc"
	"github.com/satvikmpatil/simplebank/pb"
	"github.com/satvikmpatil/simplebank/token"
	"github.com/satvikmpatil/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		panic(err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}


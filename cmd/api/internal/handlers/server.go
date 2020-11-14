package handlers

import (
	keygen "github.com/igomonov88/nimbler_key_generator/proto"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	DB *sqlx.DB
	KeyGen keygen.KeyGeneratorClient
}

func NewServer(db *sqlx.DB, keyGenClient keygen.KeyGeneratorClient) *Server{
	return &Server{DB: db, KeyGen: keyGenClient}
}

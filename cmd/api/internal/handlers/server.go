package handlers

import "github.com/jmoiron/sqlx"

type Server struct {
	DB *sqlx.DB
}

func NewServer(db *sqlx.DB) *Server{
	return &Server{DB: db}
}

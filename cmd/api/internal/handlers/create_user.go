package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	_ "google.golang.org/grpc/status"

	"nimbler_writer/internal/storage"
	pb "nimbler_writer/proto"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (resp *pb.CreateUserResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.CreateUser")
	defer span.End()

	u, err := storage.CreateUser(ctx, s.DB, req.GetName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return &pb.CreateUserResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pb.CreateUserResponse{UserID: u.ID}, nil
}

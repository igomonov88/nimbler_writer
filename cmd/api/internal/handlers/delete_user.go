package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"nimbler_writer/internal/storage"
	pb "nimbler_writer/proto"
)

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (resp *pb.DeleteUserResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.DeleteUser")
	defer span.End()

	if err := storage.DeleteUser(ctx, s.DB, req.GetUserID()); err != nil {
		return &pb.DeleteUserResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pb.DeleteUserResponse{}, nil

}

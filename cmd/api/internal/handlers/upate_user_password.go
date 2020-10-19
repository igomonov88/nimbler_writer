package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"nimbler_writer/internal/storage"
	pb "nimbler_writer/proto"
)

func (s *Server) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (resp *pb.UpdateUserPasswordResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.UpdateUserPassword")
	defer span.End()

	if err := storage.UpdateUsersPassword(ctx, s.DB, req.GetUserID(), req.GetPassword()); err != nil {
		return &pb.UpdateUserPasswordResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pb.UpdateUserPasswordResponse{}, nil

}

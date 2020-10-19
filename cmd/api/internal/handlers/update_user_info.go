package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"nimbler_writer/internal/storage"
	pb "nimbler_writer/proto"
)

func (s *Server) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (resp *pb.UpdateUserInfoResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.UpdateUserInfo")
	defer span.End()

	if err := storage.UpdateUserInfo(ctx, s.DB, req.GetUserID(), req.GetName(), req.GetEmail()); err != nil {
		return &pb.UpdateUserInfoResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pb.UpdateUserInfoResponse{}, nil
}

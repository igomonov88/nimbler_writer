package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"nimbler_writer/internal/storage"
	pb "nimbler_writer/proto"
)

func (s *Server) DoesEmailExist(ctx context.Context, req *pb.DoesEmailExistRequest) (resp *pb.DoesEmailExistResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.DoesEmailExist")
	defer span.End()

	exist, err := storage.DoesUserEmailExist(ctx, s.DB, req.GetEmail())
	if err != nil {
		return &pb.DoesEmailExistResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pb.DoesEmailExistResponse{Exist: exist}, nil
}

package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"github.com/igomonov88/nimbler_writer/internal/storage"
	pb "github.com/igomonov88/nimbler_writer/proto"
)

func (s *Server) DeleteURLList(ctx context.Context, req *pb.DeleteURLListRequest) (resp *pb.DeleteURLListResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.DeleteURLS")
	defer span.End()

	if err := storage.DeleteURLS(ctx, s.DB, req.GetUrls()); err != nil {
		return &pb.DeleteURLListResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pb.DeleteURLListResponse{}, nil
}

package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"github.com/igomonov88/nimbler_writer/internal/storage"
	pb "github.com/igomonov88/nimbler_writer/proto"
)

func (s *Server) DeleteURL(ctx context.Context, req *pb.DeleteURLRequest) (resp *pb.DeleteURLResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.DeleteURL")
	defer span.End()

	if err := storage.DeleteURL(ctx, s.DB, req.GetUrlHash()); err != nil {
		return &pb.DeleteURLResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pb.DeleteURLResponse{}, nil
}

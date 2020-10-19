package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	_ "google.golang.org/grpc/status"

	"github.com/igomonov88/nimbler_writer/internal/platform/database"
	pb "github.com/igomonov88/nimbler_writer/proto"
)

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	ctx, span := trace.StartSpan(ctx, "handlers.CheckHealth")
	defer span.End()

	if err := database.StatusCheck(ctx, s.DB); err != nil {
		return &pb.HealthCheckResponse{}, status.Error(http.StatusInternalServerError, "database is not ready")
	}

	return &pb.HealthCheckResponse{Version: "develop"}, nil
}

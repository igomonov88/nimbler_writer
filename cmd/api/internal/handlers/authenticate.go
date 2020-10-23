package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	_ "google.golang.org/grpc/status"

	"github.com/igomonov88/nimbler_writer/internal/storage"
	pb "github.com/igomonov88/nimbler_writer/proto"
)

func (s *Server) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (resp *pb.AuthenticateResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.Authenticate")
	defer span.End()

	u, err := storage.Authenticate(ctx, s.DB, req.GetEmail(), req.GetPassword())
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			return &pb.AuthenticateResponse{}, status.Error(http.StatusNotFound, err.Error())
		case storage.ErrAuthenticationFailure:
			return &pb.AuthenticateResponse{}, status.Error(http.StatusForbidden, err.Error())
		default:
			return &pb.AuthenticateResponse{}, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &pb.AuthenticateResponse{UserID: u.ID}, nil
}

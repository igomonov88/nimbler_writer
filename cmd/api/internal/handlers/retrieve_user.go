package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"github.com/igomonov88/nimbler_writer/internal/storage"
	pb "github.com/igomonov88/nimbler_writer/proto"
)

func (s *Server) RetrieveUser(ctx context.Context, req *pb.RetrieveUserRequest) (resp *pb.RetrieveUserResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.RetrieveUser")
	defer span.End()

	ru, err := storage.RetrieveUser(ctx, s.DB, req.GetUserID())
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			return &pb.RetrieveUserResponse{}, status.Error(http.StatusNotFound, err.Error())
		case storage.ErrInvalidUserID:
			return &pb.RetrieveUserResponse{}, status.Error(http.StatusBadRequest, err.Error())
		default:
			return &pb.RetrieveUserResponse{}, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &pb.RetrieveUserResponse{
		UserID: ru.ID,
		Name:   ru.Name,
		Email:  ru.Email,
	}, nil
}

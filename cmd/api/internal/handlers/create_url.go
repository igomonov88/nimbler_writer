package handlers

import (
	"context"
	"net/http"
	"time"

	keygen "github.com/igomonov88/nimbler_key_generator/proto"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"github.com/igomonov88/nimbler_writer/internal/storage"
	pb "github.com/igomonov88/nimbler_writer/proto"
)

func (s *Server) CreateUrl(ctx context.Context, req *pb.CreateURLRequest) (resp *pb.CreateURLResponse, err error) {
	ctx, span := trace.StartSpan(ctx, "handlers.CreateURL")
	defer span.End()

	kr, err := s.KeyGen.GetKey(ctx, &keygen.GetKeyRequest{})
	if err != nil {
		return &pb.CreateURLResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	u := storage.Url{
		URLHash:     kr.GetKey(),
		UserID:      req.GetUserID(),
		CreatedAt:   time.Now(),
		ExpiredAt:   req.GetExpiredAt().AsTime(),
		OriginalURL: req.GetOriginalURL(),
		CustomAlias: req.GetCustomAlias(),
	}

	if err := storage.StoreURL(ctx, s.DB, u); err != nil {
		switch err {
		case storage.ErrInvalidUserID:
			return &pb.CreateURLResponse{}, status.Error(http.StatusBadRequest, err.Error())
		default:
			return &pb.CreateURLResponse{}, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &pb.CreateURLResponse{UrlHash: kr.GetKey()}, nil
}

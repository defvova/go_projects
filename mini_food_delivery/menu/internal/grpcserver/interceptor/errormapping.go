package interceptor

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMappingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, status.Error(codes.NotFound, "not found")
	default:
		return nil, status.Error(codes.Internal, "internal error")
	}
}

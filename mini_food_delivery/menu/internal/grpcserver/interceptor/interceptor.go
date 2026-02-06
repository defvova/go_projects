package interceptor

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TimeoutInterceptor(d time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		ctx, cancel := context.WithTimeout(ctx, d)
		defer cancel()
		return handler(ctx, req)
	}
}

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

func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)

	st, _ := status.FromError(err)

	log.Info().Msgf(
		"grpc method=%s duration=%s status=%s",
		info.FullMethod,
		duration,
		st.Code(),
	)
	return resp, err
}

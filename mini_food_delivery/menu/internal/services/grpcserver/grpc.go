package grpcserver

import (
	"context"
	"database/sql"
	"errors"
	menuv1 "mini_food_delivery/menu/pkg/menu/v1"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type gRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{
		addr: addr,
	}
}

func (s *gRPCServer) ServeListener(handler menuv1.MenuServiceServer) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to listed: %v", s.addr)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingInterceptor,
			TimeoutInterceptor(2*time.Second),
			ErrorMappingInterceptor,
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Hour,
			Timeout:           20 * time.Second,
		}),
		grpc.MaxRecvMsgSize(4<<20), // 4mb
	)
	menuv1.RegisterMenuServiceServer(server, handler)

	return server, lis, nil
}

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

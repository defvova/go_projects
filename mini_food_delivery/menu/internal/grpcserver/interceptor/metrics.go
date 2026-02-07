package interceptor

import (
	"context"
	"mini_food_delivery/menu/internal/grpcserver/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MetricsInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	code := status.Code(err).String()

	metrics.GrpcRequests.
		WithLabelValues(info.FullMethod, code).
		Inc()

	metrics.GrpcLatency.
		WithLabelValues(info.FullMethod).
		Observe(time.Since(start).Seconds())

	return resp, err
}

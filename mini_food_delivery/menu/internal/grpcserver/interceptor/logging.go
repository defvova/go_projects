package interceptor

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

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

	respSize := 0
	if resp != nil {
		if msg, ok := resp.(proto.Message); ok {
			respSize = proto.Size(msg)
		}
	}

	log.Info().Msgf(
		"grpc method=%s duration=%s status=%s resp_bytes=%d",
		info.FullMethod,
		duration,
		st.Code(),
		respSize,
	)
	return resp, err
}

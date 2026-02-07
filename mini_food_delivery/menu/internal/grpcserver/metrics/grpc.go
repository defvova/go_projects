package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GrpcRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GrpcLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "gRPC request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func Register() {
	prometheus.MustRegister(GrpcRequests, GrpcLatency)
}

package grpcserver

import (
	"mini_food_delivery/menu/db"
	"mini_food_delivery/menu/internal/grpcserver/interceptor"
	"mini_food_delivery/menu/internal/services/category"
	"mini_food_delivery/menu/internal/services/menu"
	"mini_food_delivery/menu/internal/services/menuitem"
	"mini_food_delivery/menu/internal/services/menuitemprice"
	categoryv1 "mini_food_delivery/menu/pkg/category/v1"
	menuv1 "mini_food_delivery/menu/pkg/menu/v1"
	menuitemv1 "mini_food_delivery/menu/pkg/menuitem/v1"
	menuitempricev1 "mini_food_delivery/menu/pkg/menuitemprice/v1"
	"net"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GRPCOptions struct {
	OtelEnabled       bool
	PrometheusEnabled bool
}

type GRPCServer struct {
	Server *grpc.Server
}

func NewGRPCServer(q *db.Queries, opts GRPCOptions) *GRPCServer {
	var unaryInterceptors = []grpc.UnaryServerInterceptor{
		interceptor.LoggingInterceptor,
		interceptor.TimeoutInterceptor(2 * time.Second),
		interceptor.ErrorMappingInterceptor,
	}
	var (
		streamInterceptors []grpc.StreamServerInterceptor
		serverOpts         []grpc.ServerOption
	)

	if opts.OtelEnabled {
		serverOpts = append(
			serverOpts,
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		)
	}

	if opts.PrometheusEnabled {
		grpcprom.EnableHandlingTimeHistogram()
		unaryInterceptors = append(
			unaryInterceptors,
			grpcprom.UnaryServerInterceptor,
			interceptor.MetricsInterceptor,
		)
		streamInterceptors = append(
			streamInterceptors,
			grpcprom.StreamServerInterceptor,
		)
	}
	serverOpts = append(
		serverOpts,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Hour,
			Timeout:           20 * time.Second,
		}),
		grpc.MaxRecvMsgSize(4<<20),
	)

	server := grpc.NewServer(serverOpts...)

	if opts.PrometheusEnabled {
		grpcprom.Register(server)
	}

	menuHandler := menu.NewHandler(menu.NewStore(q))
	menuv1.RegisterMenuServiceServer(server, menuHandler)

	categoryHandler := category.NewHandler(category.NewStore(q))
	categoryv1.RegisterCategoryServiceServer(server, categoryHandler)

	menuItemHandler := menuitem.NewHandler(menuitem.NewStore(q))
	menuitemv1.RegisterMenuItemServiceServer(server, menuItemHandler)

	menuItemPriceHandler := menuitemprice.NewHandler(menuitemprice.NewStore(q))
	menuitempricev1.RegisterMenuItemPriceServiceServer(server, menuItemPriceHandler)

	return &GRPCServer{Server: server}
}

func (g *GRPCServer) Listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func (g *GRPCServer) Serve(lis net.Listener) error {
	return g.Server.Serve(lis)
}

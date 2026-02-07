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
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type gRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{
		addr: addr,
	}
}

func (s *gRPCServer) ServeListener(q *db.Queries) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to listed: %v", s.addr)
	}

	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			grpcprom.UnaryServerInterceptor,
			interceptor.LoggingInterceptor,
			interceptor.TimeoutInterceptor(2*time.Second),
			interceptor.ErrorMappingInterceptor,
		),
		grpc.ChainStreamInterceptor(
			grpcprom.StreamServerInterceptor,
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Hour,
			Timeout:           20 * time.Second,
		}),
		grpc.MaxRecvMsgSize(4<<20), // 4mb
	)

	grpcprom.Register(server)

	menuHandler := menu.NewHandler(menu.NewStore(q))
	menuv1.RegisterMenuServiceServer(server, menuHandler)

	categoryHandler := category.NewHandler(category.NewStore(q))
	categoryv1.RegisterCategoryServiceServer(server, categoryHandler)

	menuItemHandler := menuitem.NewHandler(menuitem.NewStore(q))
	menuitemv1.RegisterMenuItemServiceServer(server, menuItemHandler)

	menuItemPriceHandler := menuitemprice.NewHandler(menuitemprice.NewStore(q))
	menuitempricev1.RegisterMenuItemPriceServiceServer(server, menuItemPriceHandler)

	return server, lis, nil
}

package menuclient

import (
	menuv1 "mini_food_delivery/menu/pkg/menu/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MenuClient struct {
	Client menuv1.MenuServiceClient
	conn   *grpc.ClientConn
}

func NewMenuClient(addr string) (*MenuClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(4<<20)),
	)

	if err != nil {
		return nil, err
	}

	return &MenuClient{Client: menuv1.NewMenuServiceClient(conn), conn: conn}, nil
}

func (m *MenuClient) Close() error {
	return m.conn.Close()
}

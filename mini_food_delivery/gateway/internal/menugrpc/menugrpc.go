package menugrpc

import (
	categoryv1 "mini_food_delivery/menu/pkg/category/v1"
	menuv1 "mini_food_delivery/menu/pkg/menu/v1"
	menuitemv1 "mini_food_delivery/menu/pkg/menuitem/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Menu     menuv1.MenuServiceClient
	Category categoryv1.CategoryServiceClient
	MenuItem menuitemv1.MenuItemServiceClient
	conn     *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(4<<20),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		Menu:     menuv1.NewMenuServiceClient(conn),
		Category: categoryv1.NewCategoryServiceClient(conn),
		MenuItem: menuitemv1.NewMenuItemServiceClient(conn),
		conn:     conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

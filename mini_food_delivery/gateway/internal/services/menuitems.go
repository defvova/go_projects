package services

import (
	"context"
	"time"

	"github.com/defvova/go_projects/mini_food_delivery/gateway/graph/model"
	menuitemv1 "github.com/defvova/go_projects/mini_food_delivery/menu/pkg/menuitem/v1"
)

type MenuItemService struct {
	MenuItemClient menuitemv1.MenuItemServiceClient
}

func (s *MenuItemService) GetAllMenuItemsWithPrice(ctx context.Context, categoryId int64) ([]*model.MenuItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	resp, err := s.MenuItemClient.GetAllMenuItemsWithPrice(ctx, &menuitemv1.GetAllMenuItemsWithPriceRequest{CategoryId: categoryId})

	if err != nil {
		return nil, err
	}

	data := make([]*model.MenuItem, len(resp.Items))
	for i, item := range resp.Items {
		data[i] = &model.MenuItem{
			ID:         item.Id,
			CategoryID: item.CategoryId,
			Name:       item.Name,
			CreatedAt:  item.CreatedAt.AsTime(),
		}
	}

	return data, nil
}

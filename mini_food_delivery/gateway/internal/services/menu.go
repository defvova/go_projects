package services

import (
	"context"
	"time"

	menuv1 "github.com/defvova/go_projects/mini_food_delivery/menu/pkg/menu/v1"

	"github.com/defvova/go_projects/mini_food_delivery/gateway/graph/model"
)

type MenuService struct {
	MenuClient menuv1.MenuServiceClient
}

func (s *MenuService) GetAllMenus(ctx context.Context) ([]*model.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	resp, err := s.MenuClient.GetAllMenus(ctx, &menuv1.GetAllMenusRequest{})

	if err != nil {
		return nil, err
	}

	data := make([]*model.Menu, len(resp.Items))
	for i, item := range resp.Items {
		data[i] = &model.Menu{
			ID:          item.Id,
			Name:        item.Name,
			Description: &item.Description,
			Active:      item.Active,
			CreatedAt:   item.CreatedAt.AsTime(),
		}
	}

	return data, nil
}

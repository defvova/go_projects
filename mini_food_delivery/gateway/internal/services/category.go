package services

import (
	"context"
	"time"

	categoryv1 "github.com/defvova/go_projects/mini_food_delivery/menu/pkg/category/v1"

	"github.com/defvova/go_projects/mini_food_delivery/gateway/graph/model"
)

type CategoryService struct {
	CategoryClient categoryv1.CategoryServiceClient
}

func (s *CategoryService) GetAllCategories(ctx context.Context, menuId int64) ([]*model.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	resp, err := s.CategoryClient.GetAllCategories(ctx, &categoryv1.GetAllCategoriesRequest{MenuId: menuId})

	if err != nil {
		return nil, err
	}

	data := make([]*model.Category, len(resp.Items))
	for i, item := range resp.Items {
		data[i] = &model.Category{
			ID:        item.Id,
			MenuID:    item.MenuId,
			Name:      item.Name,
			Position:  item.Position,
			CreatedAt: item.CreatedAt.AsTime(),
		}
	}

	return data, nil
}

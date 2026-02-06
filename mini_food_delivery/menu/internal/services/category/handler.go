package category

import (
	"context"
	categoryv1 "mini_food_delivery/menu/pkg/category/v1"
)

type Handler struct {
	categoryv1.UnimplementedCategoryServiceServer
	store CategoryStore
}

func NewHandler(c CategoryStore) *Handler {
	return &Handler{
		store: c,
	}
}

func (h *Handler) GetAllCategories(
	ctx context.Context,
	req *categoryv1.GetAllCategoriesRequest,
) (*categoryv1.GetAllCategoriesResponse, error) {
	items, err := h.store.GetAllCategories(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}
	data := make([]*categoryv1.Category, len(items))

	for i, item := range items {
		data[i] = &categoryv1.Category{
			Id:       item.ID,
			Name:     item.Name,
			Position: item.Position,
		}
	}

	return &categoryv1.GetAllCategoriesResponse{
		Items: data,
	}, nil
}

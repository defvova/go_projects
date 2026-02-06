package menu

import (
	"context"
	menuv1 "mini_food_delivery/menu/pkg/menu/v1"
)

type Handler struct {
	menuv1.UnimplementedMenuServiceServer
	store MenuStore
}

func NewHandler(s MenuStore) *Handler {
	return &Handler{
		store: s,
	}
}

func (h *Handler) GetAllMenus(
	ctx context.Context,
	req *menuv1.GetAllMenusRequest,
) (*menuv1.GetAllMenusResponse, error) {
	items, err := h.store.GetAllMenus(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]*menuv1.Menu, len(items))

	for i, item := range items {
		data[i] = &menuv1.Menu{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description.String,
			Active:      item.Active,
		}
	}

	return &menuv1.GetAllMenusResponse{
		Items: data,
	}, nil
}

package menuitem

import (
	"context"
	menuitemv1 "mini_food_delivery/menu/pkg/menuitem/v1"
)

type Handler struct {
	menuitemv1.UnimplementedMenuItemServiceServer
	store MenuItemStore
}

func NewHandler(c MenuItemStore) *Handler {
	return &Handler{
		store: c,
	}
}

func (h *Handler) GetAllMenuItems(
	ctx context.Context,
	req *menuitemv1.GetAllMenuItemsRequest,
) (*menuitemv1.GetAllMenuItemsResponse, error) {
	items, err := h.store.GetAllMenuItems(ctx, req.CategoryId)
	if err != nil {
		return nil, err
	}
	data := make([]*menuitemv1.MenuItem, len(items))

	for i, item := range items {
		data[i] = &menuitemv1.MenuItem{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description.String,
			ImageUrl:    item.ImageUrl.String,
			Available:   item.Available,
		}
	}

	return &menuitemv1.GetAllMenuItemsResponse{
		Items: data,
	}, nil
}

package menuitem

import (
	"context"
	menuitemv1 "mini_food_delivery/menu/pkg/menuitem/v1"

	"github.com/rs/zerolog/log"
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

func (h *Handler) GetAllMenuItemsWithPrice(
	ctx context.Context,
	req *menuitemv1.GetAllMenuItemsWithPriceRequest,
) (*menuitemv1.GetAllMenuItemsWithPriceResponse, error) {
	items, err := h.store.GetAllMenuItemsWithPrice(ctx, req.CategoryId)
	if err != nil {
		log.Error().Err(err).Msg("MenuItem: sql query error")
		return nil, err
	}
	data := make([]*menuitemv1.MenuItemWithPrice, len(items))

	for i, item := range items {
		data[i] = &menuitemv1.MenuItemWithPrice{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description.String,
			ImageUrl:    item.ImageUrl.String,
			Available:   item.Available,
			PriceCents:  item.PriceCents,
			Currency:    item.Currency,
		}
	}

	return &menuitemv1.GetAllMenuItemsWithPriceResponse{
		Items: data,
	}, nil
}

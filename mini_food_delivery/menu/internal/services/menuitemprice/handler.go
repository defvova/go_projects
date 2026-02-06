package menuitemprice

import (
	"context"
	menuitempricev1 "mini_food_delivery/menu/pkg/menuitemprice/v1"
)

type Handler struct {
	menuitempricev1.UnimplementedMenuItemPriceServiceServer
	store MenuItemPriceStore
}

func NewHandler(c MenuItemPriceStore) *Handler {
	return &Handler{
		store: c,
	}
}

func (h *Handler) GetAllMenuItemPrices(
	ctx context.Context,
	req *menuitempricev1.GetAllMenuItemPricesRequest,
) (*menuitempricev1.GetAllMenuItemPricesResponse, error) {
	items, err := h.store.GetAllMenuItemPrices(ctx, req.MenuItemId)
	if err != nil {
		return nil, err
	}
	data := make([]*menuitempricev1.MenuItemPrice, len(items))

	for i, item := range items {
		data[i] = &menuitempricev1.MenuItemPrice{
			Id:         item.ID,
			PriceCents: item.PriceCents,
			Currency:   item.Currency,
		}
	}

	return &menuitempricev1.GetAllMenuItemPricesResponse{
		Items: data,
	}, nil
}

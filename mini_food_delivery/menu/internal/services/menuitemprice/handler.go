package menuitemprice

import (
	"context"
	menuitempricev1 "mini_food_delivery/menu/pkg/menuitemprice/v1"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	tr := otel.Tracer("menuItemPrice.Handler")
	ctx, span := tr.Start(ctx, "menuItemPrice.GetAllMenuItemPrices")
	defer span.End()

	ctx, dbSpan := tr.Start(ctx, "db.select.menuItemPrices")
	items, err := h.store.GetAllMenuItemPrices(ctx, req.MenuItemId)
	if err != nil {
		dbSpan.RecordError(err)
		dbSpan.SetStatus(codes.Error, err.Error())
		dbSpan.End()
		log.Error().Err(err).Msg("MenuItemPrice: sql query error")
		return nil, err
	}
	dbSpan.SetAttributes(attribute.Int("db.rows", len(items)))
	dbSpan.End()

	_, mapSpan := tr.Start(ctx, "map.menuItemPrice.entities")
	data := make([]*menuitempricev1.MenuItemPrice, len(items))
	for i, item := range items {
		data[i] = &menuitempricev1.MenuItemPrice{
			Id:         item.ID,
			PriceCents: item.PriceCents,
			Currency:   item.Currency,
		}
	}
	mapSpan.SetAttributes(attribute.Int("items.mapped", len(data)))
	mapSpan.End()

	return &menuitempricev1.GetAllMenuItemPricesResponse{
		Items: data,
	}, nil
}

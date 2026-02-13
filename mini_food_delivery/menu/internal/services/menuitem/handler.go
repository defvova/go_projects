package menuitem

import (
	"context"

	menuitemv1 "github.com/defvova/go_projects/mini_food_delivery/menu/pkg/menuitem/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	tr := otel.Tracer("menuItem.Handler")
	ctx, span := tr.Start(ctx, "menuItem.GetAllMenuItemsWithPrice")
	defer span.End()

	ctx, dbSpan := tr.Start(ctx, "db.select.menuItems")
	items, err := h.store.GetAllMenuItemsWithPrice(ctx, req.CategoryId)
	if err != nil {
		dbSpan.RecordError(err)
		dbSpan.SetStatus(codes.Error, err.Error())
		dbSpan.End()
		log.Error().Err(err).Msg("MenuItem: sql query error")
		return nil, err
	}
	dbSpan.SetAttributes(attribute.Int("db.rows", len(items)))
	dbSpan.End()

	_, mapSpan := tr.Start(ctx, "map.menuItem.entities")
	data := make([]*menuitemv1.MenuItemWithPrice, len(items))
	for i, item := range items {
		data[i] = &menuitemv1.MenuItemWithPrice{
			Id:          item.ID,
			CategoryId:  item.CategoryID,
			Name:        item.Name,
			Description: item.Description.String,
			ImageUrl:    item.ImageUrl.String,
			Available:   item.Available,
			PriceCents:  item.PriceCents,
			Currency:    item.Currency,
			CreatedAt:   timestamppb.New(item.CreatedAt),
		}
	}
	mapSpan.SetAttributes(attribute.Int("items.mapped", len(data)))
	mapSpan.End()

	return &menuitemv1.GetAllMenuItemsWithPriceResponse{
		Items: data,
	}, nil
}

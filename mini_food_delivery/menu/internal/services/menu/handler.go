package menu

import (
	"context"
	menuv1 "mini_food_delivery/menu/pkg/menu/v1"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	tr := otel.Tracer("menu.Handler")
	ctx, span := tr.Start(ctx, "menu.GetAllMenus")
	defer span.End()

	ctx, dbSpan := tr.Start(ctx, "db.select.menus")
	items, err := h.store.GetAllMenus(ctx)
	if err != nil {
		dbSpan.RecordError(err)
		dbSpan.SetStatus(codes.Error, err.Error())
		dbSpan.End()
		log.Error().Err(err).Msg("Menu: sql query error")
		return nil, err
	}
	dbSpan.SetAttributes(attribute.Int("db.rows", len(items)))
	dbSpan.End()

	_, mapSpan := tr.Start(ctx, "map.menu.entities")
	data := make([]*menuv1.Menu, len(items))
	for i, item := range items {
		data[i] = &menuv1.Menu{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description.String,
			Active:      item.Active,
		}
	}
	mapSpan.SetAttributes(attribute.Int("items.mapped", len(data)))
	mapSpan.End()

	return &menuv1.GetAllMenusResponse{
		Items: data,
	}, nil
}

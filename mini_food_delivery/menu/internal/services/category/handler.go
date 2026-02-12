package category

import (
	"context"

	categoryv1 "github.com/defvova/go_projects/mini_food_delivery/menu/pkg/category/v1"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	tr := otel.Tracer("category.Handler")
	ctx, span := tr.Start(ctx, "category.GetAllCategories")
	defer span.End()

	ctx, dbSpan := tr.Start(ctx, "db.select.categories")
	items, err := h.store.GetAllCategories(ctx, req.MenuId)
	if err != nil {
		dbSpan.RecordError(err)
		dbSpan.SetStatus(codes.Error, err.Error())
		dbSpan.End()
		log.Error().Err(err).Msg("Category: sql query error")
		return nil, err
	}
	dbSpan.SetAttributes(attribute.Int("db.rows", len(items)))
	dbSpan.End()

	_, mapSpan := tr.Start(ctx, "map.category.entities")
	data := make([]*categoryv1.Category, len(items))
	for i, item := range items {
		data[i] = &categoryv1.Category{
			Id:        item.ID,
			Name:      item.Name,
			Position:  item.Position,
			MenuId:    item.MenuID,
			CreatedAt: timestamppb.New(item.CreatedAt),
		}
	}
	mapSpan.SetAttributes(attribute.Int("items.mapped", len(data)))
	mapSpan.End()

	return &categoryv1.GetAllCategoriesResponse{
		Items: data,
	}, nil
}

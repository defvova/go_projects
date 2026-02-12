package resolver

import (
	"context"

	"github.com/defvova/go_projects/mini_food_delivery/gateway/graph/model"
)

func (r *Resolver) Categories(ctx context.Context, menuId int64) ([]*model.Category, error) {
	return r.CategoryService.GetAllCategories(ctx, menuId)
}

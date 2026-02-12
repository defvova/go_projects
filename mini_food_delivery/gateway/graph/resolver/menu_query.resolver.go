package resolver

import (
	"context"

	"github.com/defvova/go_projects/mini_food_delivery/gateway/graph/model"
)

func (r *Resolver) Menus(ctx context.Context) ([]*model.Menu, error) {
	return r.MenuService.GetAllMenus(ctx)
}

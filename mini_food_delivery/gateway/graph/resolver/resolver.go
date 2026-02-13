package resolver

import (
	"context"

	"github.com/defvova/go_projects/mini_food_delivery/gateway/graph/model"
	"github.com/defvova/go_projects/mini_food_delivery/gateway/internal/services"
)

type Resolver struct {
	MenuService     services.MenuService
	CategoryService services.CategoryService
	MenuItemService services.MenuItemService
}

func (r *Resolver) Menus(ctx context.Context) ([]*model.Menu, error) {
	return r.MenuService.GetAllMenus(ctx)
}

func (r *Resolver) Categories(ctx context.Context, menuId int64) ([]*model.Category, error) {
	return r.CategoryService.GetAllCategories(ctx, menuId)
}

func (r *Resolver) MenuItems(ctx context.Context, categoryId int64) ([]*model.MenuItem, error) {
	return r.MenuItemService.GetAllMenuItemsWithPrice(ctx, categoryId)
}

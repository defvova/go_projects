package resolver

import (
	"github.com/defvova/go_projects/mini_food_delivery/gateway/internal/services"
)

type Resolver struct {
	MenuService     services.MenuService
	CategoryService services.CategoryService
}

package services

import (
	"mini_food_delivery/gateway/internal/menugrpc"
	"mini_food_delivery/gateway/internal/services/categories"
	"mini_food_delivery/gateway/internal/services/health"
	"mini_food_delivery/gateway/internal/services/menuitems"
	"mini_food_delivery/gateway/internal/services/menus"

	"github.com/go-chi/chi/v5"
)

type ApiService struct {
	MenuGrpc *menugrpc.Client
	R        *chi.Mux
}

func (a *ApiService) Init() {
	menuClient := menus.NewHandler(a.MenuGrpc.Menu)
	categoryClient := categories.NewHandler(a.MenuGrpc.Category)
	menuItemClient := menuitems.NewHandler(a.MenuGrpc.MenuItem)

	a.R.Route("/api/v1", func(r chi.Router) {
		r.Get("/menus", menuClient.GetAllMenus)
		r.Get("/menus/{menuId}/categories", categoryClient.GetAllCategories)
		r.Get("/categories/{categoryId}/menu_items", menuItemClient.GetAllMenuItems)

		r.Get("/health", health.HandleHealth)
	})
}

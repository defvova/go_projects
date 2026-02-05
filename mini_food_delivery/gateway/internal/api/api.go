package services

import (
	"mini_food_delivery/gateway/internal/services/health"
	"mini_food_delivery/gateway/internal/services/menus"
	"mini_food_delivery/gateway/internal/services/menus/menuclient"

	"github.com/go-chi/chi/v5"
)

type ApiService struct {
	Menu menuclient.MenuClient
	R    *chi.Mux
}

func (a *ApiService) Init() {
	menuClient := menus.NewHandler(a.Menu)

	a.R.Route("/api/v1", func(r chi.Router) {
		r.Get("/menus", menuClient.GetAllMenus)
		r.Get("/health", health.HandleHealth)
	})
}

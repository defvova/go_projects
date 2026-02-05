package menus

import (
	"context"
	"mini_food_delivery/gateway/internal/services/menus/menuclient"
	"mini_food_delivery/gateway/internal/utils/grpcstatus"
	"mini_food_delivery/gateway/internal/utils/writejson"
	menuv1 "mini_food_delivery/menu/pkg/menu/v1"
	"net/http"
	"time"
)

type Handler struct {
	menu menuclient.MenuClient
}

func NewHandler(menuClient menuclient.MenuClient) *Handler {
	return &Handler{
		menu: menuClient,
	}
}

func (h *Handler) GetAllMenus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	resp, err := h.menu.Client.GetAllMenus(ctx, &menuv1.GetAllMenusRequest{})

	if err != nil {
		grpcstatus.WriteGrpcError(w, err)
		return
	}

	writejson.NewJSON(w, http.StatusOK, resp.Items)
}

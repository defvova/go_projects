package menuitems

import (
	"context"
	"mini_food_delivery/gateway/internal/utils/grpcstatus"
	"mini_food_delivery/gateway/internal/utils/writejson"
	menuitemv1 "mini_food_delivery/menu/pkg/menuitem/v1"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	menuItemClient menuitemv1.MenuItemServiceClient
}

func NewHandler(mic menuitemv1.MenuItemServiceClient) *Handler {
	return &Handler{
		menuItemClient: mic,
	}
}

func (h *Handler) GetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	id := chi.URLParam(r, "categoryId")
	categoryId, _ := strconv.Atoi(id)
	resp, err := h.menuItemClient.GetAllMenuItemsWithPrice(ctx, &menuitemv1.GetAllMenuItemsWithPriceRequest{CategoryId: int64(categoryId)})

	if err != nil {
		grpcstatus.WriteGrpcError(w, err)
		return
	}

	writejson.NewJSON(w, http.StatusOK, resp.Items)
}

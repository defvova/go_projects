package categories

import (
	"context"
	"mini_food_delivery/gateway/internal/utils/grpcstatus"
	"mini_food_delivery/gateway/internal/utils/writejson"
	categoryv1 "mini_food_delivery/menu/pkg/category/v1"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	categoryClient categoryv1.CategoryServiceClient
}

func NewHandler(cl categoryv1.CategoryServiceClient) *Handler {
	return &Handler{
		categoryClient: cl,
	}
}

func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	id := chi.URLParam(r, "menuId")
	menuId, _ := strconv.Atoi(id)
	resp, err := h.categoryClient.GetAllCategories(ctx, &categoryv1.GetAllCategoriesRequest{MenuId: int64(menuId)})

	if err != nil {
		grpcstatus.WriteGrpcError(w, err)
		return
	}

	writejson.NewJSON(w, http.StatusOK, resp.Items)
}

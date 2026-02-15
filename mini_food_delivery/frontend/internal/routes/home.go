package routes

import (
	"net/http"
	"time"

	"github.com/defvova/go_projects/mini_food_delivery/frontend/internal/graphql"
	"github.com/defvova/go_projects/mini_food_delivery/frontend/internal/service"
	"github.com/defvova/go_projects/mini_food_delivery/gateway/graph/model"
	"github.com/rs/zerolog/log"
)

type HomePage struct {
	Title         string
	FooterYear    int
	ErrMessage    string
	IsCurrentUser bool
	Menus         []model.Menu
}

type MenusQuery struct {
	Menus []model.Menu `json:"menus"`
}

type HomeHandler struct {
	Client *graphql.GraphQLClient
}

func (h *HomeHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	var resp MenusQuery
	err := h.Client.Do(r.Context(), graphql.GraphQLRequest{
		Query: `query { menus { id name description } }`,
		Variables: map[string]any{
			"active": true,
		},
	}, &resp)
	if err != nil {
		log.Err(err).Msg("gateway is failed")
		return
	}

	pageNames := []string{"home"}
	render := service.NewRenderer(pageNames)
	page := HomePage{
		Title:         "MiniFoodDelivery | Home page",
		FooterYear:    time.Now().Year(),
		IsCurrentUser: false,
		Menus:         resp.Menus,
	}
	render.Render(w, "home", page)
}

package routes

import (
	"net/http"
	"time"

	"github.com/defvova/go_projects/mini_food_delivery/frontend/internal/service"
)

type HomePage struct {
	Title         string
	FooterYear    int
	ErrMessage    string
	IsCurrentUser bool
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	pageNames := []string{"home"}
	render := service.NewRenderer(pageNames)
	page := HomePage{
		Title:         "MiniFoodDelivery | Home page",
		FooterYear:    time.Now().Year(),
		IsCurrentUser: false,
	}
	render.Render(w, "home", page)
}

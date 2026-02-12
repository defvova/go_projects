package routes

import (
	"github.com/go-chi/chi/v5"
)

// type HomePage struct {
// 	Title         string
// 	FooterYear    int
// 	ErrMessage    string
// 	IsCurrentUser bool
// }

func UserRoutes() chi.Router {
	r := chi.NewRouter()

	// r.Get("/", homeHandler)
	return r
}

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	pageNames := []string{"home"}
// 	render := service.NewRenderer(pageNames)
// 	page := HomePage{
// 		Title:         "MiniFoodDelivery | Home page",
// 		FooterYear:    time.Now().Year(),
// 		IsCurrentUser: false,
// 	}
// 	render.Render(w, "home", page)
// }

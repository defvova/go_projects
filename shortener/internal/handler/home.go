package handler

import (
	"net/http"
	"shortener/db"
	"shortener/internal/session"
	"time"
)

type HomePage struct {
	Title         string
	FooterYear    int
	ErrMessage    string
	IsCurrentUser bool
	Redirects     []db.Redirect
}

func (h *Handler) HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	userId, IsCurrentUser := session.CurrentUserID(r)
	redirects, _ := h.Q.GetRedirects(r.Context(), userId)
	page := HomePage{
		Title:         "Shortener | Home page",
		FooterYear:    time.Now().Year(),
		IsCurrentUser: IsCurrentUser,
		Redirects:     redirects,
	}
	h.NewRender.Render(w, "home", page)
}

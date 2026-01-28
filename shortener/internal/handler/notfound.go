package handler

import (
	"net/http"
	"time"
)

type NotFoundPage struct {
	Title         string
	FooterYear    int
	ErrMessage    string
	IsCurrentUser bool
}

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	page := NotFoundPage{Title: "Shortener | 404 page", FooterYear: time.Now().Year()}
	w.WriteHeader(http.StatusNotFound)
	h.NewRender.Render(w, "404", page)
}

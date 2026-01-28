package handler

import (
	"net/http"
	"shortener/internal/session"
)

func (h *Handler) LogoutPostHandler(w http.ResponseWriter, r *http.Request) {
	s := session.SessionStore{Q: h.Q}
	s.Destroy(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

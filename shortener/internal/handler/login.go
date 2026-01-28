package handler

import (
	"net/http"
	"shortener/internal/service"
	"shortener/internal/session"
	"strings"
	"time"
)

type LoginPage struct {
	Title         string
	FooterYear    int
	ErrMessage    string
	IsCurrentUser bool
}

type LoginFormParams struct {
	Email    string
	Password string
}

func (h *Handler) LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	page := LoginPage{Title: "Shortener | Login page", FooterYear: time.Now().Year()}
	h.NewRender.Render(w, "login", page)
}

func (h *Handler) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	page := LoginPage{Title: "Shortener | Login page", FooterYear: time.Now().Year()}
	r.ParseForm()
	params := LoginFormParams{Email: strings.TrimSpace(r.FormValue("email")), Password: r.FormValue("password")}
	user, err := h.Q.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		page.ErrMessage = "Invalid email or password!"
		return
	}

	ok, err := service.VerifyPassword(params.Password, user.PasswordHash)
	if err != nil {
		page.ErrMessage = err.Error()
	}
	if !ok {
		page.ErrMessage = "Invalid email or password!"
		h.NewRender.Render(w, "login", page)
		return
	}

	s := session.SessionStore{Q: h.Q}
	s.Create(w, r, user.ID, 7*24*time.Hour)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

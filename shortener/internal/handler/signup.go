package handler

import (
	"net/http"
	"shortener/db"
	"shortener/internal/service"
	"strings"
	"time"
)

type SignupPage struct {
	Title         string
	FooterYear    int
	ErrMessage    string
	IsCurrentUser bool
}

type SignupFormParams struct {
	Email                string
	Password             string
	PasswordConfirmation string
}

func (h *Handler) SignupGetHandler(w http.ResponseWriter, r *http.Request) {
	page := SignupPage{Title: "Shortener | SignUp page", FooterYear: time.Now().Year()}
	h.NewRender.Render(w, "signup", page)
}

func (h *Handler) SignupPostHandler(w http.ResponseWriter, r *http.Request) {
	page := SignupPage{Title: "Shortener | SignUp page", FooterYear: time.Now().Year()}
	r.ParseForm()
	params := SignupFormParams{Email: strings.TrimSpace(r.FormValue("email")), Password: r.FormValue("password"), PasswordConfirmation: r.FormValue("password_confirmation")}
	if !validateFields(params, &page) {
		h.NewRender.Render(w, "signup", page)
		return
	}
	passwordHash, errPass := service.HashPassword(params.Password, service.DefaultPasswordParams)
	if errPass != nil {
		page.ErrMessage = "Failed to process password"
		h.NewRender.Render(w, "signup", page)
		return
	}

	_, err := h.Q.CreateUser(r.Context(), db.CreateUserParams{
		Email:        params.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		page.ErrMessage = err.Error()
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validateFields(params SignupFormParams, page *SignupPage) bool {
	if params.Email == "" {
		page.ErrMessage = "Email is empty!"
		return false
	}

	if !strings.Contains(params.Email, "@") {
		page.ErrMessage = "Email is empty!"
		return false
	}

	if params.Password == "" {
		page.ErrMessage = "Password is empty!"
		return false
	}

	if params.PasswordConfirmation == "" {
		page.ErrMessage = "Password confirmation is empty!"
		return false
	}

	if params.Password != params.PasswordConfirmation {
		page.ErrMessage = "Password is not confirmed!"
		return false
	}
	return true
}

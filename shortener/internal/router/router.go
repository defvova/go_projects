package router

import (
	"net/http"
	"shortener/db"
	"shortener/internal/handler"
	"shortener/internal/service"
	"shortener/internal/session"

	"github.com/jackc/pgx/v5"
)

func InitHandlers(mux *http.ServeMux, conn *pgx.Conn) http.Handler {
	pageNames := []string{"login", "signup", "home", "404"}
	r := service.NewRenderer(pageNames)
	q := db.New(conn)
	h := handler.Handler{
		Q:         q,
		NewRender: r,
	}

	mux.Handle("GET /", http.HandlerFunc(h.NotFoundHandler))
	mux.Handle("POST /shortener", session.RequireAuth(http.HandlerFunc(h.ShortenerPostHandler)))
	mux.Handle("POST /shortener/{id}", session.RequireAuth(http.HandlerFunc(h.ShortenerDeleteHandler)))
	mux.Handle("GET /home", session.RequireAuth(http.HandlerFunc(h.HomeGetHandler)))
	mux.Handle("POST /logout", session.RequireAuth(http.HandlerFunc(h.LogoutPostHandler)))
	mux.Handle("GET /login", session.RedirectIfAuth(http.HandlerFunc(h.LoginGetHandler)))
	mux.Handle("POST /login", session.RedirectIfAuth(http.HandlerFunc(h.LoginPostHandler)))
	mux.Handle("GET /signup", session.RedirectIfAuth(http.HandlerFunc(h.SignupGetHandler)))
	mux.Handle("POST /signup", session.RedirectIfAuth(http.HandlerFunc(h.SignupPostHandler)))
	mux.Handle("GET /{token}", http.HandlerFunc(h.ShortenerGetHandler))

	s := session.SessionStore{Q: q}
	return s.Authenticate(mux)
}

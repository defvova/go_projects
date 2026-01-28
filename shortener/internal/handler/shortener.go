package handler

import (
	"crypto/rand"
	"log"
	"net/http"
	"net/url"
	"shortener/db"
	"shortener/internal/session"
	"strconv"
)

func NewToken(n int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b), nil
}

func (h *Handler) ShortenerGetHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	log.Println(token)
	if token == "" {
		http.NotFound(w, r)
		return
	}

	url, err := h.Q.GetRedirectByToken(r.Context(), token)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) ShortenerPostHandler(w http.ResponseWriter, r *http.Request) {
	userId, isCurrentUser := session.CurrentUserID(r)
	if !isCurrentUser {
		return
	}
	r.ParseForm()

	formUrl := r.FormValue("url")
	u, err := url.ParseRequestURI(formUrl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	token, _ := NewToken(7)
	_, err = h.Q.CreateRedirect(r.Context(), db.CreateRedirectParams{
		Url:    formUrl,
		UserID: userId,
		Token:  token,
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (h *Handler) ShortenerDeleteHandler(w http.ResponseWriter, r *http.Request) {
	userId, isCurrentUser := session.CurrentUserID(r)
	if !isCurrentUser {
		return
	}
	r.ParseForm()
	method := r.FormValue("_method")
	if method != http.MethodDelete {
		http.Error(w, "Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	redirectId, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	h.Q.DeleteRedirect(r.Context(), db.DeleteRedirectParams{
		UserID: userId,
		ID:     redirectId,
	})
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

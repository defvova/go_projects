package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"shortener/db"
	"time"
)

type ctxKey int

const userIDKey ctxKey = 1
const CookieName string = "sid"

type SessionStore struct {
	Q *db.Queries
}

func newToken() (string, []byte, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, err
	}
	raw := base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	return raw, sum[:], nil
}

func (s *SessionStore) Create(w http.ResponseWriter, r *http.Request, userID int64, ttl time.Duration) error {
	raw, hash, err := newToken()
	if err != nil {
		return err
	}

	expires := time.Now().Add(ttl)

	s.Q.CreateSession(r.Context(), db.CreateSessionParams{
		UserID:    userID,
		ExpiresAt: expires,
		TokenHash: hash,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    raw,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
	return nil
}

func (s *SessionStore) Destroy(w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie(CookieName)
	if err == nil {
		sum := sha256.Sum256([]byte(c.Value))
		s.Q.DeleteSession(r.Context(), sum[:])
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
	return nil
}

func (s *SessionStore) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(CookieName)
		if err != nil || c.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		sum := sha256.Sum256([]byte(c.Value))

		session, err := s.Q.GetSession(r.Context(), sum[:])
		userID := session.UserID
		expires := session.ExpiresAt

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if time.Now().After(expires) {
			s.Q.DeleteSession(r.Context(), sum[:])
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CurrentUserID(r *http.Request) (int64, bool) {
	v := r.Context().Value(userIDKey)
	id, ok := v.(int64)
	return id, ok
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := CurrentUserID(r); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RedirectIfAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := CurrentUserID(r); ok {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

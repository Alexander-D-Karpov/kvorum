package middleware

import (
	"context"
	"net/http"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/Alexander-D-Karpov/kvorum/internal/security"
)

type ctxKey string

const userIDKey ctxKey = "userID"

type SessionStore interface {
	GetSession(ctx context.Context, sessionID string) (*security.Session, error)
}

type Middleware struct {
	hmacSecret   string
	sessionStore SessionStore
}

func NewMiddleware(hmacSecret string, sessionStore SessionStore) *Middleware {
	return &Middleware{
		hmacSecret:   hmacSecret,
		sessionStore: sessionStore,
	}
}

func (m *Middleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		session, err := m.sessionStore.GetSession(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserID(ctx context.Context) shared.ID {
	if userID, ok := ctx.Value(userIDKey).(shared.ID); ok {
		return userID
	}
	return ""
}

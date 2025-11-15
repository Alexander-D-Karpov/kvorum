package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/Alexander-D-Karpov/kvorum/internal/security"
)

func (h *Handlers) ExchangeDeepLinkToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	tokenData, err := security.VerifyDeepLinkToken(req.Token, []byte(h.hmacSecret))
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	user, err := h.identitySvc.GetOrCreateUser(r.Context(), "max", tokenData.UserID, "")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	session, err := security.NewSession(user.ID, 30*24*time.Hour)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	if err := h.cache.SetSession(r.Context(), session); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Path:     "/",
		MaxAge:   int(30 * 24 * time.Hour.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": user.ID,
	})
}

func (h *Handlers) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.identitySvc.GetUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		_ = h.cache.DeleteSession(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	respondJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}

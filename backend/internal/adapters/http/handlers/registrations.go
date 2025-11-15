package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/registrations"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) RegisterForEvent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID := shared.ID(chi.URLParam(r, "id"))

	var req struct {
		Source string                 `json:"source"`
		UTM    map[string]interface{} `json:"utm"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	utmBytes, _ := json.Marshal(req.UTM)

	reg, err := h.registrationsSvc.Register(r.Context(), eventID, userID, req.Source, utmBytes)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to register")
		return
	}

	respondJSON(w, http.StatusCreated, reg)
}

func (h *Handlers) UpdateRSVP(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID := shared.ID(chi.URLParam(r, "id"))

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.registrationsSvc.UpdateRSVP(r.Context(), eventID, userID, registrations.Status(req.Status)); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update rsvp")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) CancelRegistration(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID := shared.ID(chi.URLParam(r, "id"))

	if err := h.registrationsSvc.CancelRegistration(r.Context(), eventID, userID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to cancel registration")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

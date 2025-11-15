package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	event, err := h.eventsSvc.CreateEvent(r.Context(), userID, req.Title, req.Description)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create event")
		return
	}

	respondJSON(w, http.StatusCreated, event)
}

func (h *Handlers) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	event, err := h.eventsSvc.GetEvent(r.Context(), eventID)
	if err != nil {
		respondError(w, http.StatusNotFound, "event not found")
		return
	}

	respondJSON(w, http.StatusOK, event)
}

func (h *Handlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID := shared.ID(chi.URLParam(r, "id"))

	var updates events.Event
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.eventsSvc.UpdateEvent(r.Context(), userID, eventID, &updates); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update event")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) PublishEvent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID := shared.ID(chi.URLParam(r, "id"))

	if err := h.eventsSvc.PublishEvent(r.Context(), userID, eventID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to publish event")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "published"})
}

func (h *Handlers) CancelEvent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID := shared.ID(chi.URLParam(r, "id"))

	if err := h.eventsSvc.CancelEvent(r.Context(), userID, eventID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to cancel event")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *Handlers) ListEvents(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, []interface{}{})
}

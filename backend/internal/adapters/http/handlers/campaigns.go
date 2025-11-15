package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	eventID := shared.ID(chi.URLParam(r, "id"))

	var req struct {
		Name        string     `json:"name"`
		Segment     string     `json:"segment"`
		Channel     string     `json:"channel"`
		Message     string     `json:"message"`
		ScheduledAt *time.Time `json:"scheduled_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	campaign, err := h.campaignsSvc.CreateCampaign(
		r.Context(),
		eventID,
		req.Name,
		req.Segment,
		req.Channel,
		req.Message,
		req.ScheduledAt,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create campaign")
		return
	}

	respondJSON(w, http.StatusCreated, campaign)
}

func (h *Handlers) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	eventID := shared.ID(chi.URLParam(r, "id"))

	campaigns, err := h.campaignsSvc.GetCampaigns(r.Context(), eventID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get campaigns")
		return
	}

	respondJSON(w, http.StatusOK, campaigns)
}

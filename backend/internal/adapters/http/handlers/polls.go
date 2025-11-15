package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/polls"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) CreatePoll(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	var req struct {
		Question string          `json:"question"`
		Options  json.RawMessage `json:"options"`
		Type     string          `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	pollType := polls.PollType(req.Type)
	poll, err := h.pollsSvc.CreatePoll(r.Context(), eventID, req.Question, req.Options, pollType)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create poll")
		return
	}

	respondJSON(w, http.StatusCreated, poll)
}

func (h *Handlers) GetEventPolls(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	polls, err := h.pollsSvc.GetPollsByEvent(r.Context(), eventID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get polls")
		return
	}

	respondJSON(w, http.StatusOK, polls)
}

func (h *Handlers) VoteOnPoll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	pollID := shared.ID(chi.URLParam(r, "id"))

	var req struct {
		OptionKey string `json:"option_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.pollsSvc.Vote(r.Context(), pollID, userID, req.OptionKey); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to vote")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "voted"})
}

func (h *Handlers) GetPollResults(w http.ResponseWriter, r *http.Request) {
	pollID := shared.ID(chi.URLParam(r, "id"))

	results, err := h.pollsSvc.GetResults(r.Context(), pollID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get results")
		return
	}

	respondJSON(w, http.StatusOK, results)
}

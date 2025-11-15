package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/checkin"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) CreateForm(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	var req struct {
		Schema json.RawMessage `json:"schema"`
		Rules  json.RawMessage `json:"rules"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	form, err := h.formsSvc.CreateForm(r.Context(), eventID, req.Schema, req.Rules)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create form")
		return
	}

	respondJSON(w, http.StatusCreated, form)
}

func (h *Handlers) GetActiveForm(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	form, err := h.formsSvc.GetActiveForm(r.Context(), eventID)
	if err != nil {
		respondError(w, http.StatusNotFound, "form not found")
		return
	}

	respondJSON(w, http.StatusOK, form)
}

func (h *Handlers) SubmitForm(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	formID := shared.ID(chi.URLParam(r, "id"))

	var req struct {
		Answers json.RawMessage `json:"answers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	response, err := h.formsSvc.SubmitResponse(r.Context(), formID, userID, req.Answers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to submit form")
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

func (h *Handlers) GetDraft(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	formID := shared.ID(chi.URLParam(r, "id"))

	draft, ok := h.formsSvc.GetDraft(r.Context(), formID, userID)
	if !ok {
		respondJSON(w, http.StatusOK, map[string]interface{}{})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"draft": draft})
}

func (h *Handlers) SaveDraft(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	formID := shared.ID(chi.URLParam(r, "id"))

	var req struct {
		Data json.RawMessage `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.formsSvc.SaveDraft(r.Context(), formID, userID, req.Data); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save draft")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) ScanCheckin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		QRCode string `json:"qr_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	c, err := h.checkinSvc.ValidateAndCheckin(r.Context(), req.QRCode, checkin.MethodQR)
	if err != nil {
		if errors.Is(err, checkin.ErrAlreadyCheckedIn) {
			respondError(w, http.StatusConflict, "already checked in")
			return
		}
		if errors.Is(err, checkin.ErrInvalidQRToken) {
			respondError(w, http.StatusBadRequest, "invalid qr code")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to checkin")
		return
	}

	respondJSON(w, http.StatusOK, c)
}

func (h *Handlers) ManualCheckin(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	c, err := h.checkinSvc.ManualCheckin(r.Context(), eventID, shared.ID(req.UserID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to checkin")
		return
	}

	respondJSON(w, http.StatusOK, c)
}

func (h *Handlers) GetQRCode(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	eventID := shared.ID(chi.URLParam(r, "id"))

	token, err := h.checkinSvc.GenerateQRToken(r.Context(), userID, eventID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate qr")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"token": token})
}

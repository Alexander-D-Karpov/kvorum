package handlers

import (
	"net/http"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetEventICS(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	ics, err := h.calendarSvc.GenerateEventICS(r.Context(), eventID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate ics")
		return
	}

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=event.ics")
	w.Write(ics)
}

func (h *Handlers) GetUserICS(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	ics, err := h.calendarSvc.GenerateUserICS(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate ics")
		return
	}

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=calendar.ics")
	w.Write(ics)
}

func (h *Handlers) GetGoogleCalendarLink(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	link, err := h.calendarSvc.GetGoogleCalendarLink(r.Context(), eventID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate link")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"link": link})
}

package handlers

import (
	"net/http"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetEventAnalytics(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from time.Time
	var to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid from")
			return
		}
	}

	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid to")
			return
		}
	}

	analytics, err := h.analyticsSvc.GetEventAnalytics(r.Context(), eventID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get analytics")
		return
	}

	respondJSON(w, http.StatusOK, analytics)
}

func (h *Handlers) ExportEventAnalyticsCSV(w http.ResponseWriter, r *http.Request) {
	eventID := shared.ID(chi.URLParam(r, "id"))

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from time.Time
	var to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid from")
			return
		}
	}

	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid to")
			return
		}
	}

	data, err := h.analyticsSvc.ExportEventAnalyticsCSV(r.Context(), eventID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to export analytics")
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=analytics.csv")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

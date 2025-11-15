package http

import (
	"net/http"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/handlers"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	*chi.Mux
}

func NewRouter(h *handlers.Handlers, m *middleware.Middleware) *Router {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/webhook/max", h.HandleMaxWebhook)

		r.Post("/auth/max/exchange", h.ExchangeDeepLinkToken)
		r.Post("/auth/logout", m.RequireAuth(h.Logout))
		r.Get("/me", m.RequireAuth(h.GetMe))

		r.Route("/events", func(r chi.Router) {
			r.Get("/", h.ListEvents)
			r.Post("/", m.RequireAuth(h.CreateEvent))
			r.Get("/{id}", h.GetEvent)
			r.Put("/{id}", m.RequireAuth(h.UpdateEvent))
			r.Post("/{id}/publish", m.RequireAuth(h.PublishEvent))
			r.Post("/{id}/cancel", m.RequireAuth(h.CancelEvent))
			r.Post("/{id}/register", m.RequireAuth(h.RegisterForEvent))
			r.Post("/{id}/rsvp", m.RequireAuth(h.UpdateRSVP))
			r.Delete("/{id}/register", m.RequireAuth(h.CancelRegistration))

			r.Route("/{id}/forms", func(r chi.Router) {
				r.Get("/active", h.GetActiveForm)
				r.Post("/", m.RequireAuth(h.CreateForm))
			})

			r.Post("/{id}/checkin/scan", m.RequireAuth(h.ScanCheckin))
			r.Post("/{id}/checkin/manual", m.RequireAuth(h.ManualCheckin))

			r.Route("/{id}/polls", func(r chi.Router) {
				r.Get("/", h.GetEventPolls)
				r.Post("/", m.RequireAuth(h.CreatePoll))
			})

			r.Get("/{id}/ics", h.GetEventICS)
			r.Get("/{id}/google-calendar", h.GetGoogleCalendarLink)
			r.Get("/{id}/analytics", m.RequireAuth(h.GetEventAnalytics))
			r.Get("/{id}/analytics.csv", m.RequireAuth(h.ExportEventAnalyticsCSV))
		})

		r.Route("/forms", func(r chi.Router) {
			r.Post("/{id}/submit", m.RequireAuth(h.SubmitForm))
			r.Get("/{id}/draft", m.RequireAuth(h.GetDraft))
			r.Put("/{id}/draft", m.RequireAuth(h.SaveDraft))
		})

		r.Route("/tickets", func(r chi.Router) {
			r.Get("/{id}/qr", m.RequireAuth(h.GetQRCode))
		})

		r.Route("/polls", func(r chi.Router) {
			r.Post("/{id}/vote", m.RequireAuth(h.VoteOnPoll))
			r.Get("/{id}/results", h.GetPollResults)
		})

		r.Get("/me/ics", m.RequireAuth(h.GetUserICS))
	})

	return &Router{Mux: r}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/botmax"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/cache"
	httpAdapter "github.com/Alexander-D-Karpov/kvorum/internal/adapters/http"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/handlers"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/queue"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/repo"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/analytics"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/calendar"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/checkin"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/forms"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/identity"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/polls"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/registrations"
	"github.com/Alexander-D-Karpov/kvorum/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	db, err := repo.NewDB(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redisCache, err := cache.NewRedisCache(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer redisCache.Close()

	scheduler, err := queue.NewAsynqScheduler(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}
	defer scheduler.Close()

	botClient := botmax.NewClient(cfg.Bot.Token)

	userRepo := repo.NewUserRepo(db)
	eventRepo := repo.NewEventRepo(db)
	roleRepo := repo.NewRoleRepo(db)
	regRepo := repo.NewRegistrationRepo(db)
	waitlistRepo := repo.NewWaitlistRepo(db)
	formRepo := repo.NewFormRepo(db)
	responseRepo := repo.NewResponseRepo(db)
	pollRepo := repo.NewPollRepo(db)
	voteRepo := repo.NewVoteRepo(db)
	checkinRepo := repo.NewCheckinRepo(db)
	qrTokenRepo := repo.NewQRTokenRepo(db)
	calendarEventRepo := repo.NewCalendarEventRepo(db)

	identitySvc := identity.NewService(userRepo)
	eventsSvc := events.NewService(eventRepo, nil, roleRepo, scheduler, redisCache)
	formsSvc := forms.NewService(formRepo, responseRepo, redisCache)
	registrationsSvc := registrations.NewService(regRepo, waitlistRepo, eventRepo)
	checkinSvc := checkin.NewService(checkinRepo, qrTokenRepo, cfg.Security.HMACSecret)
	pollsSvc := polls.NewService(pollRepo, voteRepo)
	calendarSvc := calendar.NewService(calendarEventRepo)
	analyticsSvc := analytics.NewService()

	h := handlers.NewHandlers(
		identitySvc,
		eventsSvc,
		formsSvc,
		registrationsSvc,
		checkinSvc,
		pollsSvc,
		calendarSvc,
		analyticsSvc,
		botClient,
		redisCache,
		cfg.Security.WebhookSecret,
		cfg.Security.HMACSecret,
	)
	m := middleware.NewMiddleware(cfg.Security.HMACSecret, redisCache)
	router := httpAdapter.NewRouter(h, m)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting API server on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}

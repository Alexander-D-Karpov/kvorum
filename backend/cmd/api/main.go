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
	httpserver "github.com/Alexander-D-Karpov/kvorum/internal/adapters/http"
	httphandlers "github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/handlers"
	httpmiddleware "github.com/Alexander-D-Karpov/kvorum/internal/adapters/http/middleware"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/queue"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/repo"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/analytics"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/calendar"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/campaigns"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/checkin"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/events"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/forms"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/identity"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/polls"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/registrations"
	"github.com/Alexander-D-Karpov/kvorum/internal/config"
	"github.com/Alexander-D-Karpov/kvorum/internal/observ"
	"github.com/joho/godotenv"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logger := observ.NewLogger()
	logger.Info("Starting API server...")

	ctx := context.Background()

	db, err := repo.NewDB(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	cache, err := cache.NewRedisCache(cfg.Redis.URL)
	if err != nil {
		log.Fatal("Failed to connect to redis:", err)
	}
	defer cache.Close()

	botClient, err := botmax.NewClient(cfg.Bot.Token)
	if err != nil {
		log.Fatal("Failed to create bot client:", err)
	}

	botInfo, err := botClient.Bots.GetBot(ctx)
	if err != nil {
		log.Fatal("Failed to get bot info:", err)
	}
	log.Printf("Bot started: %s (@%s)", botInfo.Name, botInfo.Username)

	webhookURL := cfg.Server.PublicURL + "/api/v1/webhook/max"
	updateTypes := []string{
		string(schemes.TypeMessageCreated),
		string(schemes.TypeMessageCallback),
		string(schemes.TypeBotStarted),
		string(schemes.TypeBotAdded),
		string(schemes.TypeBotRemoved),
	}

	subscriptions, err := botClient.Subscriptions.GetSubscriptions(ctx)
	if err != nil {
		log.Printf("Failed to get subscriptions: %v", err)
	} else {
		for _, sub := range subscriptions.Subscriptions {
			log.Printf("Removing old subscription: %s", sub.Url)
			_, _ = botClient.Subscriptions.Unsubscribe(ctx, sub.Url)
		}
	}

	result, err := botClient.Subscriptions.Subscribe(ctx, webhookURL, updateTypes)
	if err != nil {
		log.Fatal("Failed to subscribe to webhook:", err)
	}
	if result.Success {
		log.Printf("Webhook subscribed: %s", webhookURL)
		log.Printf("Note: Webhook signature validation is disabled (SDK limitation)")
	} else {
		log.Fatal("Failed to subscribe to webhook:", result.Message)
	}

	scheduler, err := queue.NewAsynqScheduler(cfg.Redis.URL)
	if err != nil {
		log.Fatal("Failed to create scheduler:", err)
	}
	defer scheduler.Close()

	userRepo := repo.NewUserRepo(db)
	eventRepo := repo.NewEventRepo(db)
	roleRepo := repo.NewRoleRepo(db)
	seriesRepo := repo.NewSeriesRepo(db)
	formRepo := repo.NewFormRepo(db)
	responseRepo := repo.NewResponseRepo(db)
	registrationRepo := repo.NewRegistrationRepo(db)
	waitlistRepo := repo.NewWaitlistRepo(db)
	checkinRepo := repo.NewCheckinRepo(db)
	qrTokenRepo := repo.NewQRTokenRepo(db)
	pollRepo := repo.NewPollRepo(db)
	voteRepo := repo.NewVoteRepo(db)
	calendarEventRepo := repo.NewCalendarEventRepo(db)
	analyticsRepo := repo.NewAnalyticsRepo(db)
	campaignRepo := repo.NewCampaignRepo(db)

	identitySvc := identity.NewService(userRepo)
	eventsSvc := events.NewService(eventRepo, seriesRepo, roleRepo, scheduler, cache)
	formsSvc := forms.NewService(formRepo, responseRepo, cache)
	registrationsSvc := registrations.NewService(registrationRepo, waitlistRepo, eventRepo)
	checkinSvc := checkin.NewService(checkinRepo, qrTokenRepo, cfg.Security.HMACSecret)
	pollsSvc := polls.NewService(pollRepo, voteRepo)
	calendarSvc := calendar.NewService(calendarEventRepo)
	analyticsSvc := analytics.NewService(analyticsRepo)
	campaignsSvc := campaigns.NewService(campaignRepo, nil, nil)

	middleware := httpmiddleware.NewMiddleware(cfg.Security.HMACSecret, cache)

	handlers := httphandlers.NewHandlers(
		identitySvc,
		eventsSvc,
		formsSvc,
		registrationsSvc,
		checkinSvc,
		pollsSvc,
		calendarSvc,
		analyticsSvc,
		campaignsSvc,
		botClient,
		cache,
		cfg.Security.WebhookSecret,
		cfg.Security.HMACSecret,
	)

	router := httpserver.NewRouter(handlers, middleware)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server listening on port %s", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	logger.Info("Server stopped")
}

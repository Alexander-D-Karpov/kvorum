package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/botmax"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/queue"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/repo"
	"github.com/Alexander-D-Karpov/kvorum/internal/config"
	"github.com/Alexander-D-Karpov/kvorum/internal/observ"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logger := observ.NewLogger()
	logger.Info("Starting worker...")

	ctx := context.Background()

	db, err := repo.NewDB(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	botClient, err := botmax.NewClient(cfg.Bot.Token)
	if err != nil {
		log.Fatal("Failed to create bot client:", err)
	}

	botInfo, err := botClient.Bots.GetBot(ctx)
	if err != nil {
		log.Fatal("Failed to get bot info:", err)
	}
	log.Printf("Worker bot client initialized: %s (@%s)", botInfo.Name, botInfo.Username)

	server, err := queue.NewAsynqServer(cfg.Redis.URL)
	if err != nil {
		log.Fatal("Failed to create worker server:", err)
	}

	handlers := queue.NewTaskHandlers(botClient.Api, nil, nil)

	mux := asynq.NewServeMux()
	mux.HandleFunc("reminder", handlers.HandleReminder)
	mux.HandleFunc("campaign", handlers.HandleCampaign)

	go func() {
		logger.Info("Worker started")
		if err := server.Start(mux); err != nil {
			log.Fatal("Worker failed to start:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")
	server.Stop()
	logger.Info("Worker stopped")
}

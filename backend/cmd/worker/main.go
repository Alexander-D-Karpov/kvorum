package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/botmax"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/queue"
	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/repo"
	"github.com/Alexander-D-Karpov/kvorum/internal/config"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
	"github.com/hibiken/asynq"
)

type EventRepo struct {
	db *repo.DB
}

func (r *EventRepo) GetEvent(ctx context.Context, eventID shared.ID) (queue.Event, error) {
	query := `
		SELECT id, title, description, starts_at, tz, location, online_url
		FROM events
		WHERE id = $1
	`

	var event queue.Event
	err := r.db.Pool().QueryRow(ctx, query, eventID).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartsAt,
		&event.Timezone,
		&event.Location,
		&event.OnlineURL,
	)

	return event, err
}

type RegistrationRepo struct {
	db *repo.DB
}

func (r *RegistrationRepo) GetUserRegistrations(ctx context.Context, eventID shared.ID) ([]queue.Registration, error) {
	query := `
		SELECT r.user_id, CAST(ui.provider_user_id AS BIGINT) as chat_id
		FROM registrations r
		JOIN user_identities ui ON ui.user_id = r.user_id AND ui.provider = 'max'
		WHERE r.event_id = $1 AND r.status = 'going'
	`

	rows, err := r.db.Pool().Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []queue.Registration
	for rows.Next() {
		var reg queue.Registration
		if err := rows.Scan(&reg.UserID, &reg.ChatID); err != nil {
			return nil, err
		}
		result = append(result, reg)
	}

	return result, rows.Err()
}

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

	botClient := botmax.NewClient(cfg.Bot.Token)

	webhookURL := fmt.Sprintf("%s/api/v1/webhook/max", cfg.Server.PublicURL)
	log.Printf("Registering webhook: %s", webhookURL)

	if cfg.Bot.Token != "" {
		sub, err := botClient.Subscribe(ctx, webhookURL, cfg.Security.WebhookSecret)
		if err != nil {
			log.Printf("Warning: Failed to register webhook: %v", err)
		} else {
			log.Printf("Webhook registered: id=%s, active=%v", sub.ID, sub.IsActive)
		}
	}

	eventRepo := &EventRepo{db: db}
	regRepo := &RegistrationRepo{db: db}

	srv, err := queue.NewAsynqServer(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Failed to create asynq server: %v", err)
	}

	handlers := queue.NewTaskHandlers(botClient, eventRepo, regRepo)
	mux := asynq.NewServeMux()
	mux.HandleFunc("reminder", handlers.HandleReminder)
	mux.HandleFunc("campaign", handlers.HandleCampaign)

	go func() {
		log.Println("Starting worker...")
		if err := srv.Start(mux); err != nil {
			log.Fatalf("Worker error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	srv.Stop()
	fmt.Println("Worker stopped")
}

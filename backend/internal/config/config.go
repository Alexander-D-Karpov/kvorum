package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Bot      BotConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Port      string
	PublicURL string
}

type DatabaseConfig struct {
	URL             string
	MaxConns        int
	MinConns        int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type RedisConfig struct {
	URL         string
	MaxRetries  int
	DialTimeout time.Duration
}

type BotConfig struct {
	Token      string
	APIURL     string
	WebhookURL string
}

type SecurityConfig struct {
	HMACSecret    string
	WebhookSecret string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:      getEnv("SERVER_PORT", "8080"),
			PublicURL: getEnv("PUBLIC_APP_URL", "http://localhost:8080"),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://kvorum:kvorum@localhost:5432/kvorum?sslmode=disable"),
			MaxConns:        getEnvInt("DB_MAX_CONNS", 25),
			MinConns:        getEnvInt("DB_MIN_CONNS", 5),
			MaxConnLifetime: time.Hour,
			MaxConnIdleTime: 30 * time.Minute,
		},
		Redis: RedisConfig{
			URL:         getEnv("REDIS_URL", "redis://localhost:6379/0"),
			MaxRetries:  3,
			DialTimeout: 5 * time.Second,
		},
		Bot: BotConfig{
			Token:  getEnv("MAX_BOT_TOKEN", ""),
			APIURL: "https://platform-api.max.ru",
		},
		Security: SecurityConfig{
			HMACSecret:    getEnv("HMAC_SECRET", "change_this_secret_key"),
			WebhookSecret: getEnv("WEBHOOK_SECRET", ""),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

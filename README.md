# Kvorum / Кворум

**Собираем людей во время и в место**

Event management platform with RSVP, quotas, check-in, and MAX Bot integration.

## Architecture

```
/backend
  ├─ cmd/
  │   ├─ api/        # HTTP API + MAX webhook handler
  │   ├─ worker/     # Background jobs (reminders, campaigns)
  │   └─ migrator/   # Database migrations
  ├─ internal/
  │   ├─ app/        # Application services (use cases)
  │   ├─ domain/     # Domain entities and business logic
  │   ├─ adapters/   # Infrastructure implementations
  │   ├─ config/     # Configuration
  │   ├─ security/   # Auth, HMAC, tokens
  │   └─ pkg/        # Shared utilities
  └─ migrations/     # SQL migrations
```

## Core Features (MVP)

- **Events Management**: Create, edit, publish events with series (RRULE)
- **Registration & RSVP**: Going/Not Going/Maybe with quotas and waitlist
- **Forms**: JSON-based forms with conditional logic and drafts
- **Reminders**: Scheduled notifications T-24h, T-3h, T-30m
- **Check-in**: QR code scanning for attendance tracking
- **Polls & Feedback**: NPS, ratings, multi-choice polls
- **Calendar Integration**: ICS generation, Google/Apple Calendar links
- **Bot Integration**: MAX Bot webhooks, inline keyboards, deep links
- **Observability**: Structured logging, metrics collection

## Tech Stack

- **Language**: Go 1.23
- **Database**: PostgreSQL 16
- **Cache/Queue**: Redis 7
- **Queue**: asynq (Redis-based)
- **HTTP Router**: chi
- **Bot**: MAX Platform API

## Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 16
- Redis 7

### Local Development

1. Clone repository:
```bash
git clone https://github.com/Alexander-D-Karpov/kvorum
cd kvorum
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Update `.env` with your credentials:
```env
MAX_BOT_TOKEN=your_bot_token_from_masterbot
DATABASE_URL=postgres://kvorum:kvorum@localhost:5432/kvorum?sslmode=disable
REDIS_URL=redis://localhost:6379/0
```

4. Start infrastructure:
```bash
docker-compose up -d postgres redis
```

5. Run migrations:
```bash
make migrate-up
```

6. Start API server:
```bash
make run-api
```

7. Start worker (in another terminal):
```bash
make run-worker
```

### Docker

Run all services:
```bash
docker-compose up
```

## API Endpoints

### Auth
- `POST /api/v1/auth/max/exchange` - Exchange deep link token
- `GET /api/v1/me` - Get current user

### Events
- `POST /api/v1/events` - Create event
- `GET /api/v1/events/{id}` - Get event
- `PUT /api/v1/events/{id}` - Update event
- `POST /api/v1/events/{id}/publish` - Publish event
- `POST /api/v1/events/{id}/cancel` - Cancel event

### Registration
- `POST /api/v1/events/{id}/register` - Register for event
- `POST /api/v1/events/{id}/rsvp` - Update RSVP status
- `DELETE /api/v1/events/{id}/register` - Cancel registration

### Forms
- `POST /api/v1/events/{id}/forms` - Create form
- `GET /api/v1/events/{id}/forms/active` - Get active form
- `POST /api/v1/forms/{id}/submit` - Submit form response
- `GET /api/v1/forms/{id}/draft` - Get draft
- `PUT /api/v1/forms/{id}/draft` - Save draft

### Check-in
- `POST /api/v1/events/{id}/checkin/scan` - Scan QR code
- `POST /api/v1/events/{id}/checkin/manual` - Manual check-in
- `GET /api/v1/tickets/{id}/qr` - Get QR code

### Polls
- `POST /api/v1/events/{id}/polls` - Create poll
- `GET /api/v1/events/{id}/polls` - Get event polls
- `POST /api/v1/polls/{id}/vote` - Vote on poll
- `GET /api/v1/polls/{id}/results` - Get poll results

### Calendar
- `GET /api/v1/events/{id}/ics` - Get event ICS file
- `GET /api/v1/events/{id}/google-calendar` - Get Google Calendar link
- `GET /api/v1/me/ics` - Get user's calendar ICS

### Webhooks
- `POST /api/v1/webhook/max` - MAX Bot webhook

## MAX Bot Integration

### Setup Webhook

Subscribe to webhook in MAX:
```bash
curl -X POST https://platform-api.max.ru/subscriptions \
  -H "Authorization: YOUR_BOT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-domain.com/api/v1/webhook/max",
    "update_types": ["message_created", "message_callback", "bot_started"],
    "secret": "YOUR_WEBHOOK_SECRET"
  }'
```

### Deep Links

Format: `https://max.ru/YOUR_BOT?start=PAYLOAD`

Example for auth:
```go
token, _ := security.GenerateDeepLinkToken(userID, secret, 5*time.Minute)
deepLink := fmt.Sprintf("https://max.ru/kvorum_bot?start=%s", token)
```

### Callback Payload Format

```
evt:<event_id>;act:<action>;arg:<arg>
```

Examples:
- `evt:123e4567-e89b-12d3-a456-426614174000;act:rsvp;arg:going`
- `evt:123e4567-e89b-12d3-a456-426614174000;act:open;arg:`

## Configuration

All configuration is done via environment variables:

```env
DATABASE_URL=postgres://user:pass@host:port/db?sslmode=disable
REDIS_URL=redis://host:port/db

SERVER_PORT=8080
PUBLIC_APP_URL=https://your-domain.com

MAX_BOT_TOKEN=your_bot_token
HMAC_SECRET=secret_for_deep_links
WEBHOOK_SECRET=secret_for_webhook_validation

LOG_LEVEL=info
```

## Development

### Run Tests
```bash
make test
```

### Build Binaries
```bash
make build
```

### Database Migrations

Create new migration:
```bash
migrate create -ext sql -dir migrations -seq migration_name
```

Apply migrations:
```bash
make migrate-up
```

Rollback:
```bash
make migrate-down
```

## Project Structure

```
kvorum/
├── cmd/
│   ├── api/         # HTTP server entry point
│   ├── worker/      # Background worker entry point
│   └── migrator/    # Migration runner
├── internal/
│   ├── app/         # Application layer (use cases)
│   │   ├── events/
│   │   ├── forms/
│   │   ├── registrations/
│   │   ├── campaigns/
│   │   ├── checkin/
│   │   └── identity/
│   ├── domain/      # Domain layer (entities, rules)
│   │   ├── events/
│   │   ├── forms/
│   │   ├── registrations/
│   │   ├── checkin/
│   │   └── shared/
│   └── adapters/    # Infrastructure layer
│       ├── http/    # HTTP handlers, middleware
│       ├── botmax/  # MAX Bot client
│       ├── repo/    # PostgreSQL repositories
│       ├── cache/   # Redis cache
│       └── queue/   # Asynq jobs
├── migrations/      # SQL migrations
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── docker-compose.yml
```

## License

MIT

## Author

Alexander Karpov (github.com/Alexander-D-Karpov)
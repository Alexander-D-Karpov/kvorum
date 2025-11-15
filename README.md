# Kvorum / Кворум

Система управления мероприятиями с регистрацией, квотами, листом ожидания, аналитикой и интеграцией с MAX Bot. Проект позволяет создавать события, собирать регистрации, отправлять напоминания и вести учет посещаемости как через веб-интерфейс, так и через бота в MAX.

## Основные возможности

- Управление событиями:
  - создание, редактирование и публикация событий;
  - поддержка серий событий (RRULE, исключения дат);
  - разные режимы видимости (публичное, по ссылке и т. д.);
  - отмена событий.
- Регистрация и RSVP:
  - регистрация на событие с различными статусами (идет, не идет, возможно);
  - ограничение по вместимости мероприятия (capacity);
  - автоматический лист ожидания и перевод из листа ожидания, когда освобождаются места;
  - хранение источника регистрации и UTM-меток.
- Формы:
  - JSON-схемы форм, привязанных к событию;
  - хранение и редактирование черновиков ответов;
  - отправка финальных ответов.
- Напоминания и кампании:
  - автоматические напоминания перед началом события (за 24 часа, 3 часа, 30 минут);
  - отложенные кампании по участникам события.
- Чек-ин:
  - генерация QR-токенов для билетов;
  - сканирование QR-кода и ручной чек-ин;
  - учет повторного входа.
- Опросы и обратная связь:
  - создание опросов по событию;
  - разные типы опросов (одиночный выбор, множественный, рейтинги, NPS);
  - подсчет и отдача результатов по опции.
- Календарь:
  - генерация ICS-файлов для события и личного календаря пользователя;
  - генерация ссылок для добавления события в Google Calendar.
- Интеграция с MAX Bot:
  - вебхук для MAX Platform;
  - построение карточек событий и напоминаний;
  - обработка callback-кнопок и deep-link токенов;
  - обмен токена deep-link на сессию веб-приложения.
- Аналитика:
  - базовые метрики по регистрации и посещаемости события;
  - экспорт в CSV с разбивкой по источникам.

## Архитектура и структура репозитория

Монорепозиторий состоит из backend и frontend-частей, а также инфраструктурных файлов:

```text
kvorum/
├── backend/                # Backend на Go
│   ├── cmd/
│   │   ├── api/            # HTTP API + обработчик вебхуков MAX
│   │   ├── worker/         # фоновые задачи (напоминания, кампании)
│   │   └── migrator/       # запуск миграций базы данных
│   ├── internal/
│   │   ├── app/            # прикладные сервисы (use-case слой)
│   │   ├── domain/         # доменные сущности и бизнес-логика
│   │   ├── adapters/       # инфраструктурные адаптеры (HTTP, БД, очереди, бот и т. д.)
│   │   ├── config/         # загрузка конфигурации
│   │   ├── security/       # токены, сессии, подписи
│   │   └── observ/         # логирование и метрики
│   ├── migrations/         # SQL-миграции схемы БД
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── Makefile
├── frontend/               # SPA на React + TypeScript + Vite
│   ├── src/
│   ├── Dockerfile
│   ├── eslint.config.js
│   └── README.md           # README фронтенда (шаблон Vite)
├── docker-compose.yml      # Общий docker-compose для всего стека
├── Makefile                # Общий Makefile (docker up/down, миграции, тесты)
├── .env.example            # Пример корневого файла окружения для docker-compose
└── .gitignore
````

Backend реализован по многослойной архитектуре (domain → app → adapters), frontend — одностраничное приложение (SPA), общающееся с API.

## Технологический стек

### Backend

* Язык: Go 1.23
* HTTP:

    * `github.com/go-chi/chi/v5` – роутер
    * `github.com/go-chi/cors` – CORS-middleware
* База данных:

    * PostgreSQL 16
    * `github.com/jackc/pgx/v5` – драйвер и пул подключений
    * `github.com/golang-migrate/migrate/v4` – миграции
* Кэш и фоновые задачи:

    * Redis 7
    * `github.com/redis/go-redis/v9` – клиент Redis
    * `github.com/hibiken/asynq` – очередь фоновых задач на Redis
* Интеграция с MAX:

    * `github.com/max-messenger/max-bot-api-client-go` – клиент MAX Bot API
* Конфигурация и утилиты:

    * `github.com/joho/godotenv` – загрузка `.env`
    * `github.com/google/uuid` – генерация идентификаторов
    * `github.com/robfig/cron/v3` – планировщик (используется косвенно)
* Логирование и вспомогательные библиотеки:

    * стандартный `log/slog`
    * косвенно: `github.com/rs/zerolog`, `go.uber.org/atomic`, `gopkg.in/yaml.v2` и другие (см. `backend/go.mod`).

### Frontend

* Язык: TypeScript
* Фреймворк: React
* Бандлер / dev-сервер: Vite
* Линтинг и качество кода:

    * ESLint
    * `@eslint/js`
    * `eslint-plugin-react-hooks`
    * `globals`
* Сборка и продакшен:

    * сборка Vite в каталог `dist`;
    * статическая выдача через nginx (см. `frontend/Dockerfile`).

Точный список зависимостей фронтенда определяется файлом `frontend/package.json`.

### Инфраструктура

* Docker, Docker Compose
* Nginx для фронтенда
* Makefile для типичных команд разработчика
* Миграции PostgreSQL в папке `backend/migrations`.

## Зависимости

### Backend (по `backend/go.mod`)

Основные модули:

* `github.com/go-chi/chi/v5`
* `github.com/go-chi/cors`
* `github.com/golang-migrate/migrate/v4`
* `github.com/google/uuid`
* `github.com/hibiken/asynq`
* `github.com/jackc/pgx/v5`
* `github.com/joho/godotenv`
* `github.com/redis/go-redis/v9`
* `github.com/max-messenger/max-bot-api-client-go`

Вспомогательные (косвенные) зависимости:

* `github.com/caarlos0/env/v6`
* `github.com/cespare/xxhash/v2`
* `github.com/dgryski/go-rendezvous`
* `github.com/golang/protobuf`
* `github.com/hashicorp/errwrap`
* `github.com/hashicorp/go-multierror`
* семейство `github.com/jackc/*` (pgpassfile, pgservicefile, puddle)
* `github.com/lib/pq`
* `github.com/mattn/go-colorable`
* `github.com/mattn/go-isatty`
* `github.com/rs/zerolog`
* `github.com/spf13/cast`
* `go.uber.org/atomic`
* `golang.org/x/crypto`, `golang.org/x/sync`, `golang.org/x/sys`, `golang.org/x/text`, `golang.org/x/time`
* `google.golang.org/protobuf`
* `gopkg.in/yaml.v2`.

### Frontend

* React
* React DOM
* TypeScript
* Vite
* ESLint
* `@eslint/js`
* `eslint-plugin-react-hooks`
* `globals`

Полный список и версии доступны в `frontend/package.json`.

## Конфигурация окружения

### Корневой `.env` для docker-compose

Файл `.env.example` в корне содержит пример конфигурации для docker-композа:

```env
# Application
PUBLIC_APP_URL=http://localhost

# MAX Bot
MAX_BOT_TOKEN=your_bot_token_here

# Security
HMAC_SECRET=change_this_secret_key_for_deep_links
WEBHOOK_SECRET=change_this_webhook_secret

# Database (optional overrides)
POSTGRES_USER=kvorum
POSTGRES_PASSWORD=kvorum
POSTGRES_DB=kvorum
```

Шаги:

1. Скопировать файл и отредактировать значения:

   ```bash
   cp .env.example .env
   ```
2. Обязательно задать:

    * `PUBLIC_APP_URL` – публичный URL фронтенда (например, `http://localhost`);
    * `MAX_BOT_TOKEN` – токен бота в MAX;
    * `HMAC_SECRET` – секрет для deep-link токенов и QR;
    * `WEBHOOK_SECRET` – секрет подписи вебхука (на будущее).

### Backend `.env`

Для локального запуска backend без Docker используется `backend/.env.example`:

```env
DATABASE_URL=postgres://kvorum:kvorum@localhost:5432/kvorum?sslmode=disable
REDIS_URL=redis://localhost:6379/0

SERVER_PORT=8080
PUBLIC_APP_URL=https://kvorum.example.com

MAX_BOT_TOKEN=your_bot_token_here

HMAC_SECRET=change_this_secret_key_for_deep_links
WEBHOOK_SECRET=change_this_webhook_secret

LOG_LEVEL=info
```

Шаги:

```bash
cd backend
cp .env.example .env
# далее отредактировать значения при необходимости
```

Все команды `go run`/`go build` внутри `backend` автоматически подхватят переменные из `.env` через `godotenv`.

## Запуск в Docker

Для запуска всего стека (PostgreSQL, Redis, миграции, API, worker, frontend) используется корневой `docker-compose.yml`.

### Быстрый старт через Makefile

В корне проекта:

```bash
# Сборка образов
make build

# Запуск всех сервисов в фоне
make up

# Просмотр логов
make logs

# Остановка и удаление контейнеров
make down
```

`make` использует команды из корневого `Makefile`:

* `make up` – `docker-compose up -d`
* `make down` – `docker-compose down`
* `make logs` – `docker-compose logs -f`
* `make migrate-up` – запуск мигратора
* `make migrate-down` – откат миграций
* `make clean` – `docker-compose down -v` и очистка ненужных ресурсов Docker
* `make build` – `docker-compose build`
* `make test` – запуск тестов backend и frontend.

### Ручные команды Docker

Альтернативно можно вызвать docker-compose напрямую:

```bash
# Сборка всех образов
docker-compose build

# Запуск сервиса базы и Redis
docker-compose up -d postgres redis

# Применение миграций (через отдельный сервис migrator)
docker-compose up migrator

# Запуск API, worker и frontend
docker-compose up -d api worker frontend
```

После успешного запуска:

* API доступен по адресу `http://localhost:8080`;
* фронтенд (SPA) – по адресу `http://localhost` (порт 80 из контейнера `frontend` проброшен на хост).

## Локальный запуск без Docker

Этот вариант полезен при разработке, когда база и Redis поднимаются отдельно (через Docker или локально), а код backend/frontend запускается напрямую.

### 1. База данных и Redis

Проще всего поднять их в контейнерах:

```bash
# В корне проекта
docker-compose up -d postgres redis
```

Тем самым будут доступны:

* PostgreSQL: `postgres://kvorum:kvorum@localhost:5432/kvorum`;
* Redis: `redis://localhost:6379/0`.

### 2. Миграции базы

```bash
cd backend
cp .env.example .env    # если еще не сделано
go run ./cmd/migrator/main.go up
```

Или через backend-Makefile:

```bash
cd backend
make migrate-up
```

### 3. Запуск API сервера

```bash
cd backend
go run ./cmd/api/main.go
# или
make run-api
```

Сервер поднимется на порту, указанном в `SERVER_PORT` (по умолчанию 8080).

### 4. Запуск worker

Worker обрабатывает задачи очереди Asynq (напоминания, кампании):

```bash
cd backend
go run ./cmd/worker/main.go
# или
make run-worker
```

### 5. Запуск frontend

В директории `frontend`:

```bash
cd frontend

# Установка зависимостей
npm ci

# Запуск dev-сервера (стандартный сценарий для Vite)
npm run dev
```

После этого SPA будет доступна на стандартном порту Vite (обычно `http://localhost:5173`).

Для production-сборки фронтенда:

```bash
cd frontend
npm run build
```

Собранный фронтенд лежит в `frontend/dist`. В Docker-окружении он автоматически обслуживается nginx (см. `frontend/Dockerfile`).

## Примеры команд через командную строку

### Пример: регистрация на событие через API

```bash
curl -X POST http://localhost:8080/api/v1/events/{event_id}/register \
  -H "Content-Type: application/json" \
  --cookie "session=SESSION_ID" \
  -d '{
    "source": "landing",
    "utm": {
      "utm_source": "newsletter",
      "utm_campaign": "spring"
    }
  }'
```

## HTTP эндпоинты

* Аутентификация:

    * `POST /api/v1/auth/max/exchange` – обмен deep-link токена на сессию;
    * `GET /api/v1/me` – текущий пользователь;
    * `POST /api/v1/auth/logout` – выход из системы.
* События:

    * `GET /api/v1/events` – список публичных событий;
    * `POST /api/v1/events` – создание события;
    * `GET /api/v1/events/{id}` – получение события;
    * `PUT /api/v1/events/{id}` – обновление события;
    * `POST /api/v1/events/{id}/publish` – публикация;
    * `POST /api/v1/events/{id}/cancel` – отмена.
* Регистрация и RSVP:

    * `POST /api/v1/events/{id}/register` – регистрация;
    * `POST /api/v1/events/{id}/rsvp` – изменение статуса;
    * `DELETE /api/v1/events/{id}/register` – отмена регистрации.
* Формы:

    * `POST /api/v1/events/{id}/forms` – создание формы;
    * `GET /api/v1/events/{id}/forms/active` – получение активной формы;
    * `POST /api/v1/forms/{id}/submit` – отправка ответа;
    * `GET /api/v1/forms/{id}/draft` – получить черновик;
    * `PUT /api/v1/forms/{id}/draft` – сохранить черновик.
* Чек-ин:

    * `POST /api/v1/events/{id}/checkin/scan` – чек-ин по QR;
    * `POST /api/v1/events/{id}/checkin/manual` – ручной чек-ин;
    * `GET /api/v1/tickets/{id}/qr` – получение QR-токена.
* Опросы:

    * `POST /api/v1/events/{id}/polls` – создание опроса;
    * `GET /api/v1/events/{id}/polls` – список опросов события;
    * `POST /api/v1/polls/{id}/vote` – голосование;
    * `GET /api/v1/polls/{id}/results` – результаты.
* Календарь:

    * `GET /api/v1/events/{id}/ics` – ICS-файл события;
    * `GET /api/v1/events/{id}/google-calendar` – ссылка для Google Calendar;
    * `GET /api/v1/me/ics` – личный ICS пользователя.
* Аналитика:

    * `GET /api/v1/events/{id}/analytics` – агрегированная статистика;
    * `GET /api/v1/events/{id}/analytics.csv` – экспорт в CSV.
* Вебхуки:

    * `POST /api/v1/webhook/max` – обработчик вебхука MAX Bot.

## MAX Bot и вебхуки

Backend при старте:

1. Инициализирует клиента MAX по `MAX_BOT_TOKEN`.
2. Получает информацию о боте (название и username).
3. Удаляет старые подписки и подписывается на вебхук с адресом:

   ```text
   {PUBLIC_APP_URL}/api/v1/webhook/max
   ```
4. Подписка включает типы обновлений:

    * `message_created`;
    * `message_callback`;
    * `bot_started`;
    * `bot_added`;
    * `bot_removed`.

Обработчик вебхука (`POST /api/v1/webhook/max`) принимает обновления, распознает тип, создает/обновляет пользователя в системе и отвечает пользователю через бот.

Для авторизации в веб-приложении используется deep-link:

1. В боте генерируется подпись (HMAC) с помощью `HMAC_SECRET`.
2. Пользователь открывает ссылку вида:

   ```text
   https://max.ru/<bot_username>?start=<token>
   ```
3. Фронтенд или backend обменивает токен на сессию по эндпоинту:

   ```text
   POST /api/v1/auth/max/exchange
   ```
4. Backend создает HTTP-сессию, записывает cookie `session` и далее использует ее для аутентификации.

## Тесты

### Backend

```bash
cd backend
make test
# эквивалентно:
# go test -v -race -coverprofile=coverage.out ./...
```

Результат покрытия сохраняется в `backend/coverage.html`.

### Frontend

Из корня через общий Makefile:

```bash
make test
```

Команда запускает `npm test` во `frontend`. Конкретная конфигурация тестового раннера определяется в `frontend/package.json`.

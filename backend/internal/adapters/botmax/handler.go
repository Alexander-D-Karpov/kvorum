package botmax

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Alexander-D-Karpov/kvorum/internal/app/identity"
	"github.com/Alexander-D-Karpov/kvorum/internal/app/registrations"
	domainregistrations "github.com/Alexander-D-Karpov/kvorum/internal/domain/registrations"
	maxbotapi "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Handler struct {
	api              *maxbotapi.Api
	identitySvc      *identity.Service
	registrationsSvc *registrations.Service
	hmacSecret       string
}

func NewHandler(
	api *maxbotapi.Api,
	identitySvc *identity.Service,
	registrationsSvc *registrations.Service,
	hmacSecret string,
) *Handler {
	return &Handler{
		api:              api,
		identitySvc:      identitySvc,
		registrationsSvc: registrationsSvc,
		hmacSecret:       hmacSecret,
	}
}

func (h *Handler) Handle(ctx context.Context, upd schemes.UpdateInterface) error {
	switch u := upd.(type) {
	case *schemes.MessageCreatedUpdate:
		return h.handleMessageCreated(ctx, u)
	case *schemes.BotStartedUpdate:
		return h.handleBotStarted(ctx, u)
	case *schemes.MessageCallbackUpdate:
		return h.handleMessageCallback(ctx, u)
	default:
		log.Printf("Unknown update type: %s", upd.GetUpdateType())
	}
	return nil
}

func (h *Handler) handleMessageCreated(ctx context.Context, u *schemes.MessageCreatedUpdate) error {
	userIDStr := strconv.FormatInt(u.Message.Sender.UserId, 10)
	displayName := u.Message.Sender.FirstName
	if u.Message.Sender.LastName != "" {
		displayName += " " + u.Message.Sender.LastName
	}

	_, err := h.identitySvc.GetOrCreateUser(ctx, "max", userIDStr, displayName)
	if err != nil {
		return fmt.Errorf("get or create user: %w", err)
	}

	text := u.GetText()
	chatID := u.Message.Recipient.ChatId

	switch text {
	case "/start":
		return h.sendWelcome(ctx, chatID, u.Message.Sender.FirstName)
	case "/help":
		return h.sendHelp(ctx, chatID)
	default:
		log.Printf("Unhandled message: %s", text)
	}

	return nil
}

func (h *Handler) handleBotStarted(ctx context.Context, u *schemes.BotStartedUpdate) error {
	userIDStr := strconv.FormatInt(u.User.UserId, 10)
	displayName := u.User.FirstName
	if u.User.LastName != "" {
		displayName += " " + u.User.LastName
	}

	user, err := h.identitySvc.GetOrCreateUser(ctx, "max", userIDStr, displayName)
	if err != nil {
		return fmt.Errorf("get or create user: %w", err)
	}

	log.Printf("Bot started by user: %s (ID: %s)", displayName, user.ID)

	return h.sendWelcome(ctx, u.ChatId, u.User.FirstName)
}

func (h *Handler) handleMessageCallback(ctx context.Context, u *schemes.MessageCallbackUpdate) error {
	payload, err := ParseCallbackPayload(u.Callback.Payload)
	if err != nil {
		return h.answerCallback(ctx, u.Callback.CallbackID, "Ошибка обработки")
	}

	userIDStr := strconv.FormatInt(u.Callback.User.UserId, 10)
	displayName := u.Callback.User.FirstName
	if u.Callback.User.LastName != "" {
		displayName += " " + u.Callback.User.LastName
	}

	user, err := h.identitySvc.GetOrCreateUser(ctx, "max", userIDStr, displayName)
	if err != nil {
		return fmt.Errorf("get or create user: %w", err)
	}

	switch payload.Action {
	case "rsvp":
		status := domainregistrations.Status(payload.Arg)
		if err := h.registrationsSvc.UpdateRSVP(ctx, payload.EventID, user.ID, status); err != nil {
			return h.answerCallback(ctx, u.Callback.CallbackID, "Ошибка")
		}

		notifications := map[domainregistrations.Status]string{
			domainregistrations.StatusGoing:    "✅ Вы записаны",
			domainregistrations.StatusNotGoing: "❌ Отменено",
			domainregistrations.StatusMaybe:    "❓ Напомним позже",
		}

		notification := notifications[status]
		if notification == "" {
			notification = "Статус обновлён"
		}

		return h.answerCallback(ctx, u.Callback.CallbackID, notification)

	case "confirm":
		if err := h.registrationsSvc.UpdateRSVP(ctx, payload.EventID, user.ID, domainregistrations.StatusGoing); err != nil {
			return h.answerCallback(ctx, u.Callback.CallbackID, "Ошибка")
		}
		return h.answerCallback(ctx, u.Callback.CallbackID, "✅ Подтверждено")

	case "cancel":
		if err := h.registrationsSvc.CancelRegistration(ctx, payload.EventID, user.ID); err != nil {
			return h.answerCallback(ctx, u.Callback.CallbackID, "Ошибка")
		}
		return h.answerCallback(ctx, u.Callback.CallbackID, "❌ Отменено")

	default:
		return h.answerCallback(ctx, u.Callback.CallbackID, "Неизвестное действие")
	}
}

func (h *Handler) sendWelcome(ctx context.Context, chatID int64, firstName string) error {
	text := fmt.Sprintf("👋 Привет, %s!\n\n", firstName)
	text += "Я — бот Kvorum для управления событиями.\n\n"
	text += "Я помогу тебе:\n"
	text += "• Найти интересные мероприятия\n"
	text += "• Зарегистрироваться на события\n"
	text += "• Получать напоминания\n"
	text += "• Управлять своими регистрациями\n"

	kb := h.api.Messages.NewKeyboardBuilder()
	kb.AddRow().AddLink("🎫 Мои события", schemes.DEFAULT, "https://maxapp.akarpov.ru/me")

	msg := maxbotapi.NewMessage().
		SetChat(chatID).
		SetText(text).
		SetFormat("markdown").
		AddKeyboard(kb)

	_, err := h.api.Messages.Send(ctx, msg)
	return err
}

func (h *Handler) sendHelp(ctx context.Context, chatID int64) error {
	text := "Команды:\n"
	text += "/start - Начать\n"
	text += "/help - Помощь\n"

	msg := maxbotapi.NewMessage().
		SetChat(chatID).
		SetText(text)

	_, err := h.api.Messages.Send(ctx, msg)
	return err
}

func (h *Handler) answerCallback(ctx context.Context, callbackID, notification string) error {
	ans := &schemes.CallbackAnswer{
		Notification: notification,
	}
	_, err := h.api.Messages.AnswerOnCallback(ctx, callbackID, ans)
	return err
}

package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/botmax"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/registrations"
	"github.com/Alexander-D-Karpov/kvorum/internal/security"
)

func (h *Handlers) HandleMaxWebhook(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get("X-Max-Bot-Api-Secret")
	if signature != h.webhookSecret {
		log.Printf("Invalid webhook signature")
		respondError(w, http.StatusUnauthorized, "invalid signature")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read webhook body: %v", err)
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	update, err := botmax.ParseUpdate(body)
	if err != nil {
		log.Printf("Failed to parse webhook: %v", err)
		respondError(w, http.StatusBadRequest, "invalid update")
		return
	}

	log.Printf("Webhook: type=%s", update.UpdateType)

	switch update.UpdateType {
	case "message_created":
		mc, _ := update.AsMessageCreated()
		go h.handleMessageCreated(context.Background(), mc)
	case "message_callback":
		mc, _ := update.AsMessageCallback()
		go h.handleMessageCallback(context.Background(), mc)
	case "bot_started":
		bs, _ := update.AsBotStarted()
		go h.handleBotStarted(context.Background(), bs)
	case "bot_added":
		ba, _ := update.AsBotAdded()
		log.Printf("Bot added: chat=%d", ba.ChatID)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleMessageCreated(ctx context.Context, mc *botmax.MessageCreated) {
	userIDStr := strconv.FormatInt(mc.Message.Sender.UserID, 10)
	displayName := mc.Message.Sender.FirstName
	if mc.Message.Sender.LastName != "" {
		displayName += " " + mc.Message.Sender.LastName
	}

	_, err := h.identitySvc.GetOrCreateUser(ctx, "max", userIDStr, displayName)
	if err != nil {
		log.Printf("Failed to get/create user: %v", err)
		return
	}

	switch mc.Message.Body.Text {
	case "/start":
		msg := botmax.BuildWelcomeMessage(mc.Message.Sender.FirstName)
		_, _ = h.botClient.SendMessage(ctx, mc.Message.Recipient.ChatID, msg)
	case "/help":
		helpMsg := &botmax.SendMessageRequest{
			Text: "Команды:\n/start - Начать\n/help - Помощь\n/events - События",
		}
		_, _ = h.botClient.SendMessage(ctx, mc.Message.Recipient.ChatID, helpMsg)
	}
}

func (h *Handlers) handleMessageCallback(ctx context.Context, mc *botmax.MessageCallback) {
	payload, err := botmax.ParseCallbackPayload(mc.Callback.Payload)
	if err != nil {
		_ = h.botClient.AnswerCallback(ctx, mc.Callback.CallbackID, "Ошибка", nil)
		return
	}

	userIDStr := strconv.FormatInt(mc.Callback.User.UserID, 10)
	displayName := mc.Callback.User.FirstName
	if mc.Callback.User.LastName != "" {
		displayName += " " + mc.Callback.User.LastName
	}

	user, err := h.identitySvc.GetOrCreateUser(ctx, "max", userIDStr, displayName)
	if err != nil {
		_ = h.botClient.AnswerCallback(ctx, mc.Callback.CallbackID, "Ошибка", nil)
		return
	}

	switch payload.Action {
	case "rsvp":
		status := registrations.Status(payload.Arg)
		err := h.registrationsSvc.UpdateRSVP(ctx, payload.EventID, user.ID, status)
		if err != nil {
			_ = h.botClient.AnswerCallback(ctx, mc.Callback.CallbackID, "Ошибка", nil)
			return
		}

		event, _ := h.eventsSvc.GetEvent(ctx, payload.EventID)
		if mc.Message != nil && event != nil {
			card := botmax.BuildEventCard(event, status)
			_ = h.botClient.EditMessage(ctx, mc.Message.Body.Mid, card)
		}

		notifications := map[registrations.Status]string{
			registrations.StatusGoing:    "✅ Вы записаны",
			registrations.StatusNotGoing: "❌ Отменено",
			registrations.StatusMaybe:    "❓ Напомним позже",
		}
		_ = h.botClient.AnswerCallback(ctx, mc.Callback.CallbackID, notifications[status], nil)

	case "confirm":
		_ = h.registrationsSvc.UpdateRSVP(ctx, payload.EventID, user.ID, registrations.StatusGoing)
		_ = h.botClient.AnswerCallback(ctx, mc.Callback.CallbackID, "✅ Подтверждено", nil)

	case "cancel":
		_ = h.registrationsSvc.CancelRegistration(ctx, payload.EventID, user.ID)
		_ = h.botClient.AnswerCallback(ctx, mc.Callback.CallbackID, "❌ Отменено", nil)
	}
}

func (h *Handlers) handleBotStarted(ctx context.Context, bs *botmax.BotStarted) {
	userIDStr := strconv.FormatInt(bs.User.UserID, 10)
	displayName := bs.User.FirstName
	if bs.User.LastName != "" {
		displayName += " " + bs.User.LastName
	}

	user, err := h.identitySvc.GetOrCreateUser(ctx, "max", userIDStr, displayName)
	if err != nil {
		return
	}

	if bs.Payload != "" {
		_, err := security.VerifyDeepLinkToken(bs.Payload, []byte(h.hmacSecret))
		if err == nil {
			log.Printf("User %s authenticated via deep link", user.ID)
		}
	}

	msg := botmax.BuildWelcomeMessage(bs.User.FirstName)
	_, _ = h.botClient.SendMessage(ctx, bs.ChatID, msg)
}

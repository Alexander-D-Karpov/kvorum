package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Alexander-D-Karpov/kvorum/internal/adapters/botmax"
	"github.com/Alexander-D-Karpov/kvorum/internal/domain/registrations"
	maxbotapi "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func (h *Handlers) HandleMaxWebhook(w http.ResponseWriter, r *http.Request) {
	// Log all headers for debugging
	signature := r.Header.Get("X-Max-Bot-Api-Secret")
	log.Printf("Webhook received. Signature header: '%s'", signature)

	// Temporarily allow all requests since SDK doesn't support setting webhook secret
	// TODO: Implement custom subscription with secret support
	// if signature != h.webhookSecret {
	// 	log.Printf("Invalid webhook signature")
	// 	respondError(w, http.StatusUnauthorized, "invalid signature")
	// 	return
	// }

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read webhook body: %v", err)
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	log.Printf("Webhook body: %s", string(body))

	var baseUpdate schemes.Update
	if err := json.Unmarshal(body, &baseUpdate); err != nil {
		log.Printf("Failed to parse webhook: %v", err)
		respondError(w, http.StatusBadRequest, "invalid update")
		return
	}

	log.Printf("Webhook: type=%s", baseUpdate.UpdateType)

	switch baseUpdate.UpdateType {
	case schemes.TypeMessageCreated:
		var mc schemes.MessageCreatedUpdate
		if err := json.Unmarshal(body, &mc); err != nil {
			respondError(w, http.StatusBadRequest, "invalid message_created")
			return
		}
		go h.handleMessageCreated(context.Background(), &mc)
	case schemes.TypeMessageCallback:
		var mc schemes.MessageCallbackUpdate
		if err := json.Unmarshal(body, &mc); err != nil {
			respondError(w, http.StatusBadRequest, "invalid message_callback")
			return
		}
		go h.handleMessageCallback(context.Background(), &mc)
	case schemes.TypeBotStarted:
		var bs schemes.BotStartedUpdate
		if err := json.Unmarshal(body, &bs); err != nil {
			respondError(w, http.StatusBadRequest, "invalid bot_started")
			return
		}
		go h.handleBotStarted(context.Background(), &bs)
	case schemes.TypeBotAdded:
		var ba schemes.BotAddedToChatUpdate
		if err := json.Unmarshal(body, &ba); err == nil {
			log.Printf("Bot added: chat=%d", ba.ChatId)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleMessageCreated(ctx context.Context, mc *schemes.MessageCreatedUpdate) {
	userIDStr := strconv.FormatInt(mc.Message.Sender.UserId, 10)
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
		components := botmax.BuildWelcomeMessageComponents(h.botClient.Api, mc.Message.Sender.FirstName)
		msg := maxbotapi.NewMessage().
			SetChat(mc.Message.Recipient.ChatId).
			SetText(components.Text).
			SetFormat("markdown").
			AddKeyboard(components.Keyboard)
		_, _ = h.botClient.Messages.Send(ctx, msg)
	case "/help":
		helpMsg := maxbotapi.NewMessage().
			SetChat(mc.Message.Recipient.ChatId).
			SetText("Команды:\n/start - Начать\n/help - Помощь\n/events - События")
		_, _ = h.botClient.Messages.Send(ctx, helpMsg)
	}
}

func (h *Handlers) handleMessageCallback(ctx context.Context, mc *schemes.MessageCallbackUpdate) {
	payload, err := botmax.ParseCallbackPayload(mc.Callback.Payload)
	if err != nil {
		_, _ = h.botClient.Messages.AnswerOnCallback(ctx, mc.Callback.CallbackID, &schemes.CallbackAnswer{
			Notification: "Ошибка",
		})
		return
	}

	userIDStr := strconv.FormatInt(mc.Callback.User.UserId, 10)
	displayName := mc.Callback.User.FirstName
	if mc.Callback.User.LastName != "" {
		displayName += " " + mc.Callback.User.LastName
	}

	user, err := h.identitySvc.GetOrCreateUser(ctx, "max", userIDStr, displayName)
	if err != nil {
		_, _ = h.botClient.Messages.AnswerOnCallback(ctx, mc.Callback.CallbackID, &schemes.CallbackAnswer{
			Notification: "Ошибка",
		})
		return
	}

	switch payload.Action {
	case "rsvp":
		status := registrations.Status(payload.Arg)
		err := h.registrationsSvc.UpdateRSVP(ctx, payload.EventID, user.ID, status)
		if err != nil {
			_, _ = h.botClient.Messages.AnswerOnCallback(ctx, mc.Callback.CallbackID, &schemes.CallbackAnswer{
				Notification: "Ошибка",
			})
			return
		}

		event, _ := h.eventsSvc.GetEvent(ctx, payload.EventID)
		if mc.Message != nil && event != nil {
			components := botmax.BuildEventCardComponents(h.botClient.Api, &botmax.EventForCard{
				ID:          event.ID,
				Title:       event.Title,
				Description: event.Description,
				StartsAt:    event.StartsAt,
				Timezone:    event.Timezone,
				Location:    event.Location,
				OnlineURL:   event.OnlineURL,
			}, status)

			editMsg := maxbotapi.NewMessage().
				SetText(components.Text).
				SetFormat("markdown").
				AddKeyboard(components.Keyboard)

			msgID, parseErr := strconv.ParseInt(mc.Message.Body.Mid, 10, 64)
			if parseErr == nil {
				_ = h.botClient.Messages.EditMessage(ctx, msgID, editMsg)
			} else {
				log.Printf("Failed to parse message ID: %v", parseErr)
			}
		}

		notifications := map[registrations.Status]string{
			registrations.StatusGoing:    "✅ Вы записаны",
			registrations.StatusNotGoing: "❌ Отменено",
			registrations.StatusMaybe:    "❓ Напомним позже",
		}
		_, _ = h.botClient.Messages.AnswerOnCallback(ctx, mc.Callback.CallbackID, &schemes.CallbackAnswer{
			Notification: notifications[status],
		})

	case "confirm":
		_ = h.registrationsSvc.UpdateRSVP(ctx, payload.EventID, user.ID, registrations.StatusGoing)
		_, _ = h.botClient.Messages.AnswerOnCallback(ctx, mc.Callback.CallbackID, &schemes.CallbackAnswer{
			Notification: "✅ Подтверждено",
		})

	case "cancel":
		_ = h.registrationsSvc.CancelRegistration(ctx, payload.EventID, user.ID)
		_, _ = h.botClient.Messages.AnswerOnCallback(ctx, mc.Callback.CallbackID, &schemes.CallbackAnswer{
			Notification: "❌ Отменено",
		})
	}
}

func (h *Handlers) handleBotStarted(ctx context.Context, bs *schemes.BotStartedUpdate) {
	userIDStr := strconv.FormatInt(bs.User.UserId, 10)
	displayName := bs.User.FirstName
	if bs.User.LastName != "" {
		displayName += " " + bs.User.LastName
	}

	user, err := h.identitySvc.GetOrCreateUser(ctx, "max", userIDStr, displayName)
	if err != nil {
		return
	}

	log.Printf("Bot started by user %s (ID: %s)", displayName, user.ID)

	components := botmax.BuildWelcomeMessageComponents(h.botClient.Api, bs.User.FirstName)
	msg := maxbotapi.NewMessage().
		SetChat(bs.ChatId).
		SetText(components.Text).
		SetFormat("markdown").
		AddKeyboard(components.Keyboard)
	_, _ = h.botClient.Messages.Send(ctx, msg)
}

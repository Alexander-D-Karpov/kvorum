package botmax

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://platform-api.max.ru"

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}
}

type SendMessageRequest struct {
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Format      string       `json:"format,omitempty"` // "markdown" or "html"
	Link        *MessageLink `json:"link,omitempty"`
	Notify      bool         `json:"notify"`
}

type Attachment struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

type MessageLink struct {
	Type string `json:"type"` // "reply" or "forward"
	Mid  string `json:"mid"`
}

type InlineKeyboard struct {
	Buttons [][]Button `json:"buttons"`
}

type Button struct {
	Type    string `json:"type"` // "callback", "link", "open_app", "message", "request_contact", "request_geo_location"
	Text    string `json:"text"`
	Payload string `json:"payload,omitempty"` // For callback and message types
	URL     string `json:"url,omitempty"`     // For link type
}

type SendMessageResponse struct {
	Message struct {
		Body struct {
			Mid string `json:"mid"`
		} `json:"body"`
	} `json:"message"`
}

// SendMessage отправляет сообщение в чат
func (c *Client) SendMessage(ctx context.Context, chatID int64, req *SendMessageRequest) (string, error) {
	url := fmt.Sprintf("%s/messages?chat_id=%d", BaseURL, chatID)

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", c.token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("max api error: %d %s", resp.StatusCode, string(respBody))
	}

	var result SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.Message.Body.Mid, nil
}

// EditMessage редактирует существующее сообщение
func (c *Client) EditMessage(ctx context.Context, messageID string, req *SendMessageRequest) error {
	url := fmt.Sprintf("%s/messages?message_id=%s", BaseURL, messageID)

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", c.token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("max api error: %d %s", resp.StatusCode, string(respBody))
	}

	return nil
}

type AnswerCallbackRequest struct {
	Notification string              `json:"notification,omitempty"`
	Message      *SendMessageRequest `json:"message,omitempty"`
}

// AnswerCallback отвечает на callback от inline-кнопки
func (c *Client) AnswerCallback(ctx context.Context, callbackID string, notification string, message *SendMessageRequest) error {
	url := fmt.Sprintf("%s/answers?callback_id=%s", BaseURL, callbackID)

	req := AnswerCallbackRequest{
		Notification: notification,
		Message:      message,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", c.token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("max api error: %d %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// DeleteMessage удаляет сообщение
func (c *Client) DeleteMessage(ctx context.Context, messageID string) error {
	url := fmt.Sprintf("%s/messages?message_id=%s", BaseURL, messageID)

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("max api error: %d %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// GetMe возвращает информацию о боте
func (c *Client) GetMe(ctx context.Context) (*BotInfo, error) {
	url := fmt.Sprintf("%s/me", BaseURL)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("max api error: %d %s", resp.StatusCode, string(respBody))
	}

	var info BotInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &info, nil
}

type BotInfo struct {
	UserID           int64  `json:"user_id"`
	Name             string `json:"name"`
	Username         string `json:"username"`
	IsBot            bool   `json:"is_bot"`
	LastActivityTime int64  `json:"last_activity_time"`
}

package botmax

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Subscription struct {
	URL         string   `json:"url"`
	UpdateTypes []string `json:"update_types"`
	Secret      string   `json:"secret"`
}

type SubscriptionResponse struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	UpdateTypes []string `json:"update_types"`
	IsActive    bool     `json:"is_active"`
}

func (c *Client) Subscribe(ctx context.Context, webhookURL, secret string) (*SubscriptionResponse, error) {
	url := fmt.Sprintf("%s/subscriptions", BaseURL)

	sub := Subscription{
		URL: webhookURL,
		UpdateTypes: []string{
			"message_created",
			"message_callback",
			"bot_started",
			"bot_added",
		},
		Secret: secret,
	}

	body, err := json.Marshal(sub)
	if err != nil {
		return nil, fmt.Errorf("marshal subscription: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", c.token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("subscription failed: status %d", resp.StatusCode)
	}

	var result SubscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

func (c *Client) GetSubscriptions(ctx context.Context) ([]SubscriptionResponse, error) {
	url := fmt.Sprintf("%s/subscriptions", BaseURL)

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
		return nil, fmt.Errorf("get subscriptions failed: status %d", resp.StatusCode)
	}

	var result []SubscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

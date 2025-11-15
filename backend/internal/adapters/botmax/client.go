package botmax

import (
	maxbotapi "github.com/max-messenger/max-bot-api-client-go"
)

type Client struct {
	*maxbotapi.Api
}

func NewClient(token string) (*Client, error) {
	api, err := maxbotapi.New(token)
	if err != nil {
		return nil, err
	}

	return &Client{Api: api}, nil
}

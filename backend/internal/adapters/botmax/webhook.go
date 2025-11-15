package botmax

import (
	"encoding/json"
	"time"
)

type Update struct {
	UpdateType string          `json:"update_type"`
	Timestamp  int64           `json:"timestamp"`
	RawData    json.RawMessage `json:"-"`
}

type MessageCreated struct {
	UpdateType string `json:"update_type"`
	Timestamp  int64  `json:"timestamp"`
	Message    struct {
		Sender struct {
			UserID    int64  `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
		} `json:"sender"`
		Recipient struct {
			ChatID   int64  `json:"chat_id"`
			ChatType string `json:"chat_type"` // "dialog", "chat", "channel"
		} `json:"recipient"`
		Body struct {
			Mid  string `json:"mid"`
			Text string `json:"text"`
		} `json:"body"`
	} `json:"message"`
	UserLocale string `json:"user_locale"`
}

type MessageCallback struct {
	UpdateType string `json:"update_type"`
	Timestamp  int64  `json:"timestamp"`
	Callback   struct {
		Timestamp  int64  `json:"timestamp"`
		CallbackID string `json:"callback_id"`
		Payload    string `json:"payload"`
		User       struct {
			UserID    int64  `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
		} `json:"user"`
	} `json:"callback"`
	Message *struct {
		Body struct {
			Mid string `json:"mid"`
		} `json:"body"`
		Recipient struct {
			ChatID int64 `json:"chat_id"`
		} `json:"recipient"`
	} `json:"message"`
}

type BotStarted struct {
	UpdateType string `json:"update_type"`
	Timestamp  int64  `json:"timestamp"`
	ChatID     int64  `json:"chat_id"`
	Payload    string `json:"payload"` // Deep link payload
	User       struct {
		UserID    int64  `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	} `json:"user"`
	UserLocale string `json:"user_locale"`
}

type BotAdded struct {
	UpdateType string `json:"update_type"`
	Timestamp  int64  `json:"timestamp"`
	ChatID     int64  `json:"chat_id"`
	User       struct {
		UserID    int64  `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	} `json:"user"`
}

func ParseUpdate(data []byte) (*Update, error) {
	var u Update
	u.RawData = data
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (u *Update) Time() time.Time {
	return time.UnixMilli(u.Timestamp)
}

func (u *Update) AsMessageCreated() (*MessageCreated, error) {
	var mc MessageCreated
	if err := json.Unmarshal(u.RawData, &mc); err != nil {
		return nil, err
	}
	return &mc, nil
}

func (u *Update) AsMessageCallback() (*MessageCallback, error) {
	var mc MessageCallback
	if err := json.Unmarshal(u.RawData, &mc); err != nil {
		return nil, err
	}
	return &mc, nil
}

func (u *Update) AsBotStarted() (*BotStarted, error) {
	var bs BotStarted
	if err := json.Unmarshal(u.RawData, &bs); err != nil {
		return nil, err
	}
	return &bs, nil
}

func (u *Update) AsBotAdded() (*BotAdded, error) {
	var ba BotAdded
	if err := json.Unmarshal(u.RawData, &ba); err != nil {
		return nil, err
	}
	return &ba, nil
}

package botmax

import (
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type KeyboardBuilder struct {
	rows []*KeyboardRow
}

type KeyboardRow struct {
	cols []schemes.ButtonInterface
}

func NewKeyboardBuilder() *KeyboardBuilder {
	return &KeyboardBuilder{
		rows: []*KeyboardRow{},
	}
}

func (kb *KeyboardBuilder) AddRow() *KeyboardRow {
	row := &KeyboardRow{cols: []schemes.ButtonInterface{}}
	kb.rows = append(kb.rows, row)
	return row
}

func (kb *KeyboardBuilder) Build() *schemes.Keyboard {
	buttons := make([][]schemes.ButtonInterface, len(kb.rows))
	for i, row := range kb.rows {
		buttons[i] = row.cols
	}
	return &schemes.Keyboard{
		Buttons: buttons,
	}
}

func (r *KeyboardRow) AddCallback(text string, intent schemes.Intent, payload string) *KeyboardRow {
	btn := schemes.CallbackButton{
		Button: schemes.Button{
			Text: text,
			Type: schemes.CALLBACK,
		},
		Intent:  intent,
		Payload: payload,
	}
	r.cols = append(r.cols, btn)
	return r
}

func (r *KeyboardRow) AddLink(text string, intent schemes.Intent, url string) *KeyboardRow {
	btn := schemes.LinkButton{
		Button: schemes.Button{
			Text: text,
			Type: schemes.LINK,
		},
		Url: url,
	}
	r.cols = append(r.cols, btn)
	return r
}

func (r *KeyboardRow) AddContact(text string) *KeyboardRow {
	btn := schemes.RequestContactButton{
		Button: schemes.Button{
			Text: text,
			Type: schemes.CONTACT,
		},
	}
	r.cols = append(r.cols, btn)
	return r
}

func (r *KeyboardRow) AddGeolocation(text string, quick bool) *KeyboardRow {
	btn := schemes.RequestGeoLocationButton{
		Button: schemes.Button{
			Text: text,
			Type: schemes.GEOLOCATION,
		},
		Quick: quick,
	}
	r.cols = append(r.cols, btn)
	return r
}

type OpenAppButton struct {
	schemes.Button
	WebApp  string `json:"web_app,omitempty"`
	Payload string `json:"payload,omitempty"`
}

func (b OpenAppButton) GetType() schemes.ButtonType {
	return "open_app"
}

func (r *KeyboardRow) AddOpenApp(text, webApp, payload string) *KeyboardRow {
	btn := OpenAppButton{
		Button: schemes.Button{
			Text: text,
			Type: "open_app",
		},
		WebApp:  webApp,
		Payload: payload,
	}
	r.cols = append(r.cols, btn)
	return r
}

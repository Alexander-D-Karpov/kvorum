package botmax

import (
	"fmt"
	"strings"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type CallbackPayload struct {
	EventID shared.ID
	Action  string
	Arg     string
}

func ParseCallbackPayload(payload string) (*CallbackPayload, error) {
	parts := strings.Split(payload, ";")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid payload format")
	}

	cp := &CallbackPayload{}

	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "evt":
			cp.EventID = shared.ID(kv[1])
		case "act":
			cp.Action = kv[1]
		case "arg":
			cp.Arg = kv[1]
		}
	}

	if cp.EventID == "" || cp.Action == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	return cp, nil
}

func FormatCallbackPayload(eventID shared.ID, action, arg string) string {
	payload := fmt.Sprintf("evt:%s;act:%s", eventID, action)
	if arg != "" {
		payload += fmt.Sprintf(";arg:%s", arg)
	}
	return payload
}

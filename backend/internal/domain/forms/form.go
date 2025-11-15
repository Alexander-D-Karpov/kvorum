package forms

import (
	"encoding/json"
	"errors"

	"github.com/Alexander-D-Karpov/kvorum/internal/domain/shared"
)

type Form struct {
	ID      shared.ID
	EventID shared.ID
	Version int
	Schema  json.RawMessage
	Rules   json.RawMessage
	Active  bool
	shared.Timestamp
}

type Response struct {
	ID      shared.ID
	FormID  shared.ID
	UserID  shared.ID
	Status  ResponseStatus
	Answers json.RawMessage
	shared.Timestamp
}

type ResponseStatus string

const (
	ResponseStatusDraft     ResponseStatus = "draft"
	ResponseStatusSubmitted ResponseStatus = "submitted"
)

var (
	ErrFormNotFound     = errors.New("form not found")
	ErrResponseNotFound = errors.New("response not found")
	ErrInvalidSchema    = errors.New("invalid form schema")
)

func NewForm(eventID shared.ID, schema, rules json.RawMessage) *Form {
	return &Form{
		ID:        shared.NewID(),
		EventID:   eventID,
		Version:   1,
		Schema:    schema,
		Rules:     rules,
		Active:    true,
		Timestamp: shared.NewTimestamp(),
	}
}

func NewResponse(formID, userID shared.ID) *Response {
	return &Response{
		ID:        shared.NewID(),
		FormID:    formID,
		UserID:    userID,
		Status:    ResponseStatusDraft,
		Answers:   json.RawMessage("{}"),
		Timestamp: shared.NewTimestamp(),
	}
}

func (r *Response) Submit() {
	r.Status = ResponseStatusSubmitted
	r.Timestamp.Touch()
}

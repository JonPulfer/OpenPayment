package openPayment

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           string              `json:"id"`
	Type         string              `json:"type"`
	Acknowledged bool                `json:"acknowledged"`
	Received     time.Time           `json:"received"`
	MetaData     map[string][]string `json:"metaData"`
	Data         json.RawMessage     `json:"Data"`
}

func NewEvent(eventType string, data []byte) Event {
	eventId, _ := uuid.NewRandom()

	return Event{
		ID:           eventId.String(),
		Type:         eventType,
		Acknowledged: false,
		Received:     time.Now(),
		MetaData:     nil,
		Data:         data,
	}
}

// Account events
const (
	AccountAddEvent    = "account add"
	AccountUpdateEvent = "account update"
	AccountDeleteEvent = "account delete"
)

package eventStream

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"
)

// EventStream where we are sourcing our application state from.
type EventStream interface {
	Publish(msg openPayment.Event) error
	Listen(receive chan openPayment.Event) error
}

type InMemoryStream struct {
	stream    []openPayment.Event
	nextIndex int
}

func NewInMemoryStream() *InMemoryStream {
	return &InMemoryStream{
		stream: make([]openPayment.Event, 0),
	}
}

func (ims *InMemoryStream) Publish(event openPayment.Event) error {
	ims.stream = append(ims.stream, event)
	log.Debug().Fields(map[string]interface{}{
		"eventId":   event.ID,
		"eventType": event.Type,
		"streamLen": len(ims.stream),
	}).Msg("event published")

	return nil
}

func (ims *InMemoryStream) Listen(receive chan openPayment.Event) error {
	defer close(receive)

	for {
		time.Sleep(time.Second)
		event, err := ims.next()
		if _, emptyStream := err.(streamEmpty); emptyStream {
			continue
		}

		log.Debug().Fields(map[string]interface{}{
			"eventId":   event.ID,
			"eventType": event.Type,
		}).Msg("event received")
		receive <- event
	}
}

type streamEmpty string

func (se streamEmpty) Error() string {
	return string(se)
}

func (ims *InMemoryStream) next() (openPayment.Event, error) {
	if len(ims.stream) == 0 {
		return openPayment.Event{}, streamEmpty("nothing in stream")
	}

	switch {
	case ims.nextIndex == 0:
		ims.nextIndex++
		return ims.stream[0], nil
	case ims.nextIndex >= len(ims.stream):
		oldestUnAcknowledgedIndex := ims.findOldestUnacknowledgedIndex()
		if oldestUnAcknowledgedIndex > 0 && oldestUnAcknowledgedIndex < len(ims.stream) {
			return ims.stream[oldestUnAcknowledgedIndex], nil
		}
		return openPayment.Event{}, streamEmpty("no available events")
	default:
		idxToReturn := ims.nextIndex
		ims.nextIndex++
		return ims.stream[idxToReturn], nil
	}
}

func (ims InMemoryStream) findOldestUnacknowledgedIndex() int {
	for idx, ev := range ims.stream {
		if !ev.Acknowledged {
			return idx
		}
	}
	return len(ims.stream)
}

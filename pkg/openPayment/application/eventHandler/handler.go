package eventHandler

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/infrastructure/eventStream"
)

type Subscriber interface {
	Receive(newEvents, processedEvents chan openPayment.Event) error
	fmt.Stringer
}

type subscriberChannels struct {
	send, receive chan openPayment.Event
}

type Handler struct {
	subscribers map[Subscriber]subscriberChannels
	eventStream eventStream.EventStream
}

func NewHandler(stream eventStream.EventStream) *Handler {
	return &Handler{
		subscribers: make(map[Subscriber]subscriberChannels),
		eventStream: stream,
	}
}

// Subscribe to Handler to receive events from it.
func (h *Handler) Subscribe(subscriber Subscriber) {
	send := make(chan openPayment.Event)
	receive := make(chan openPayment.Event)
	h.subscribers[subscriber] = subscriberChannels{
		send:    send,
		receive: receive,
	}

	go func(subscriber Subscriber, newEvents, processedEvents chan openPayment.Event) {
		if err := subscriber.Receive(send, receive); err != nil {
			log.Error().Err(err).Msgf("error from %s", subscriber)
		}
	}(subscriber, send, receive)
}

// Handle each event by sending it to all subscribers.
func (h *Handler) Handle(wg *sync.WaitGroup) error {
	defer wg.Done()
	received := make(chan openPayment.Event)

	go func(received chan openPayment.Event) {
		if err := h.eventStream.Listen(received); err != nil {
			log.Error().Err(err).Msg("event stream error")
		}
	}(received)

	for newEvent := range received {
		h.deliverEvent(newEvent)
	}

	return nil
}

func (h *Handler) deliverEvent(newEvent openPayment.Event) {
	log.Debug().Fields(map[string]interface{}{
		"eventId": newEvent.ID,
	}).Msg("received event")

	var waitForSubscribers sync.WaitGroup
	for subscriber := range h.subscribers {
		pendingSubscriber := subscriber
		waitForSubscribers.Add(1)
		go h.sendToReceiver(newEvent, pendingSubscriber, &waitForSubscribers)
	}

	waitForSubscribers.Wait()
	log.Debug().Fields(map[string]interface{}{
		"eventId": newEvent.ID,
	}).Msg("event processed by all subscribers")
}

func (h *Handler) sendToReceiver(newEvent openPayment.Event,
	pendingSubscriber Subscriber, wg *sync.WaitGroup,
) {
	defer wg.Done()

	ev := newEvent
	h.subscribers[pendingSubscriber].send <- ev
	evDone := <-h.subscribers[pendingSubscriber].receive

	log.Debug().Fields(map[string]interface{}{
		"eventId": evDone.ID,
	}).Msgf("event processed by %s", pendingSubscriber)
}

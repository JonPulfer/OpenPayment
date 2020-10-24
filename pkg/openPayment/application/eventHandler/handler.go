package eventHandler

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/infrastructure/eventStream"
)

type Receiver interface {
	Receive(newEvents, processedEvents chan openPayment.Event) error
	fmt.Stringer
}

type receiverChannels struct {
	send, receive chan openPayment.Event
}

type Handler struct {
	receivers   map[Receiver]receiverChannels
	eventStream eventStream.EventStream
}

func NewHandler(stream eventStream.EventStream) *Handler {
	return &Handler{
		receivers:   make(map[Receiver]receiverChannels),
		eventStream: stream,
	}
}

func (h *Handler) AddReceiver(receiver Receiver) {
	send := make(chan openPayment.Event)
	receive := make(chan openPayment.Event)
	h.receivers[receiver] = receiverChannels{
		send:    send,
		receive: receive,
	}

	go func(receiver Receiver, newEvents, processedEvents chan openPayment.Event) {
		if err := receiver.Receive(send, receive); err != nil {
			log.Error().Err(err).Msgf("error from %s", receiver)
		}
	}(receiver, send, receive)
}

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

	var waitForReceivers sync.WaitGroup
	for receiver := range h.receivers {
		pendingReceiver := receiver
		waitForReceivers.Add(1)
		go h.sendToReceiver(newEvent, pendingReceiver, &waitForReceivers)
	}

	log.Debug().Fields(map[string]interface{}{
		"eventId": newEvent.ID,
	}).Msg("event processed by all receivers")

	waitForReceivers.Wait()
}

func (h *Handler) sendToReceiver(newEvent openPayment.Event,
	pendingReceiver Receiver, wg *sync.WaitGroup,
) {
	defer wg.Done()

	ev := newEvent
	h.receivers[pendingReceiver].send <- ev
	evDone := <-h.receivers[pendingReceiver].receive

	log.Debug().Fields(map[string]interface{}{
		"eventId": evDone.ID,
	}).Msgf("event processed by %s", pendingReceiver)
}

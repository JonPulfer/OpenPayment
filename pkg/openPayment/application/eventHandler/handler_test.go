package eventHandler

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/infrastructure/eventStream"
)

func TestHandler_Handle(t *testing.T) {
	ims := eventStream.NewInMemoryStream()
	handler := NewHandler(ims)
	testReceiver := TestReceiver{}
	handler.AddReceiver(testReceiver)

	var handlerWg sync.WaitGroup
	handlerWg.Add(1)
	go func() {
		err := handler.Handle(&handlerWg)
		if err != nil {
			panic(err)
		}
	}()

	testPayload := struct {
		Name string `json:"name"`
	}{
		"tester",
	}
	testData, err := json.Marshal(&testPayload)
	require.Nil(t, err)

	event := openPayment.NewEvent("test", testData)
	require.Nil(t, ims.Publish(event))
	time.Sleep(2 * time.Second)
}

type TestReceiver struct{}

func (tr TestReceiver) Receive(newEvents, processedEvents chan openPayment.Event) error {
	for received := range newEvents {
		log.Debug().Fields(map[string]interface{}{
			"eventId": received.ID,
		}).Msgf("event received by %s", tr)
		processedEvents <- received
	}
	return nil
}

func (tr TestReceiver) String() string {
	return "TestReceiver"
}

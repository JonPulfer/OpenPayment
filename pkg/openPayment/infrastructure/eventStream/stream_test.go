package eventStream

import (
	"encoding/json"
	"testing"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStream_Publish(t *testing.T) {
	ims := NewInMemoryStream()

	testPayload := struct {
		Name string `json:"name"`
	}{
		"tester",
	}
	testData, err := json.Marshal(&testPayload)
	require.Nil(t, err)

	event := openPayment.NewEvent("test", testData)

	err = ims.Publish(event)
	assert.Nil(t, err)
}

func TestInMemoryStream_Listen(t *testing.T) {
	ims := NewInMemoryStream()
	receive := make(chan openPayment.Event)
	go ims.Listen(receive)

	testPayload := struct {
		Name string `json:"name"`
	}{
		"tester",
	}
	testData, err := json.Marshal(&testPayload)
	require.Nil(t, err)

	event := openPayment.NewEvent("test", testData)

	err = ims.Publish(event)
	select {
	case receivedEvent := <-receive:
		assert.Equal(t, event.ID, receivedEvent.ID)
	}

}

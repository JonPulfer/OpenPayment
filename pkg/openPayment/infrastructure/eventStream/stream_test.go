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

	require.Nil(t, ims.Publish(event))
	receivedEvent := <-receive
	assert.Equal(t, event.ID, receivedEvent.ID)

	testPayload2 := struct {
		Name string `json:"name"`
	}{
		"tester2",
	}
	testData2, err := json.Marshal(&testPayload2)
	require.Nil(t, err)

	event2 := openPayment.NewEvent("test", testData2)

	require.Nil(t, ims.Publish(event2))

	testPayload3 := struct {
		Name string `json:"name"`
	}{
		"tester3",
	}
	testData3, err := json.Marshal(&testPayload3)
	require.Nil(t, err)

	event3 := openPayment.NewEvent("test", testData3)

	require.Nil(t, ims.Publish(event3))

	received2 := <-receive
	assert.Equal(t, event2.ID, received2.ID)
}

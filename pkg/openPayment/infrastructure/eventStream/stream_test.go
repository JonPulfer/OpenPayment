package eventStream

import (
	"encoding/json"
	"testing"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStream_Publish(t *testing.T) {

	testPayload := struct {
		Name string `json:"name"`
	}{
		"tester",
	}
	testData, err := json.Marshal(&testPayload)
	require.Nil(t, err)

	event := openPayment.NewEvent("test", testData)

	ims := NewInMemoryStream()

	err = ims.Publish(event)
	assert.Nil(t, err)
}

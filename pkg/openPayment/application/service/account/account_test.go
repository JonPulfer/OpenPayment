package account

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/application/eventHandler"
	"github.com/JonPulfer/OpenPayment/pkg/openPayment/infrastructure/eventStream"
)

func TestSimpleAccountAddAccount(t *testing.T) {
	ims := eventStream.NewInMemoryStream()
	handler := eventHandler.NewHandler(ims)
	simpleAccount := NewSimpleAccount()
	handler.Subscribe(simpleAccount)
	var handlerWg sync.WaitGroup
	handlerWg.Add(1)
	go func() {
		err := handler.Handle(&handlerWg)
		if err != nil {
			panic(err)
		}
	}()

	newAccount := openPayment.Account{
		CustomerURLs: []string{"http://customer-service/customer/01"},
		CardURLs:     []string{"http://card-service/card/1234-1234-1324-1242"},
		Number:       123456789,
		SortCode:     [3]int{10, 20, 30},
	}
	testData, err := json.Marshal(&newAccount)
	require.Nil(t, err)

	event := openPayment.NewEvent(openPayment.AccountAddEvent, testData)
	require.Nil(t, ims.Publish(event))
	time.Sleep(2 * time.Second)
	storedAccount, err := simpleAccount.Fetch(123456789)
	require.Nil(t, err)
	assert.Equal(t, newAccount.CardURLs[0], storedAccount.CardURLs[0])
	assert.Equal(t, newAccount.CustomerURLs[0], storedAccount.CustomerURLs[0])
	assert.Equal(t, newAccount.SortCode, storedAccount.SortCode)
}

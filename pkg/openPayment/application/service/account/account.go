package account

import (
	"bytes"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"
)

type Account interface {
	Fetch(accountNumber int) (*openPayment.Account, error)
}

type SimpleAccount struct {
	accounts map[int]*openPayment.Account
}

func NewSimpleAccount() *SimpleAccount {
	return &SimpleAccount{accounts: make(map[int]*openPayment.Account)}
}

func (sa SimpleAccount) Fetch(accountNumber int) (*openPayment.Account, error) {
	if acc, ok := sa.accounts[accountNumber]; ok {
		return acc, nil
	}
	return nil, Error("Account does not exist")
}

type Error string

func (e Error) Error() string {
	return string(e)
}

func (sa SimpleAccount) String() string {
	return "simple account"
}

func (sa SimpleAccount) Receive(newEvents, processedEvents chan openPayment.Event) error {
	defer close(processedEvents)

	for ev := range newEvents {
		if !accountEvent(ev) {
			log.Debug().Fields(map[string]interface{}{
				"eventId":   ev.ID,
				"eventType": ev.Type,
			}).Msg("ignoring event")
			processedEvents <- ev
			continue
		}

		log.Debug().Fields(map[string]interface{}{
			"eventId":   ev.ID,
			"eventType": ev.Type,
		}).Msg("received account event")

		switch ev.Type {
		case openPayment.AccountAddEvent:
			err := sa.handleAccountAddEvent(ev)
			if err != nil {
				log.Error().Fields(map[string]interface{}{
					"eventId": ev.ID,
				}).Err(err).Msg("failed to add account")
				continue
			}
			processedEvents <- ev
		case openPayment.AccountUpdateEvent:
			err := sa.handleAccountUpdateEvent(ev)
			if err != nil {
				log.Error().Fields(map[string]interface{}{
					"eventId": ev.ID,
				}).Err(err).Msg("failed to update account")
				continue
			}
		case openPayment.AccountDeleteEvent:
			err := sa.handleAccountDeleteEvent(ev)
			if err != nil {
				log.Error().Fields(map[string]interface{}{
					"eventId": ev.ID,
				}).Err(err).Msg("failed to delete account")
				continue
			}
		}

		log.Debug().Fields(map[string]interface{}{
			"eventId":   ev.ID,
			"eventType": ev.Type,
		}).Msg("processed event")
	}

	return nil
}

func (sa SimpleAccount) handleAccountAddEvent(ev openPayment.Event) error {

	acc, err := extractAccount(ev.Data)
	if err != nil {
		return err
	}

	sa.accounts[acc.Number] = acc

	return nil
}

func (sa SimpleAccount) handleAccountUpdateEvent(ev openPayment.Event) error {

	acc, err := extractAccount(ev.Data)
	if err != nil {
		return err
	}

	sa.accounts[acc.Number] = acc

	return nil
}

func (sa SimpleAccount) handleAccountDeleteEvent(ev openPayment.Event) error {

	acc, err := extractAccount(ev.Data)
	if err != nil {
		return err
	}

	delete(sa.accounts, acc.Number)

	return nil
}

func accountEvent(ev openPayment.Event) bool {
	switch ev.Type {
	case openPayment.AccountAddEvent, openPayment.AccountUpdateEvent, openPayment.AccountDeleteEvent:
		return true
	}
	return false
}

func extractAccount(data json.RawMessage) (*openPayment.Account, error) {
	var acc openPayment.Account
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&acc); err != nil {
		return nil, err
	}

	return &acc, nil
}

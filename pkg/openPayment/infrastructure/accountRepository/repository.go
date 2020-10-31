package accountRepository

import "github.com/JonPulfer/OpenPayment/pkg/openPayment"

type AccountRepository interface {
	Fetch(accountNumber int) (*openPayment.Account, error)
	Store(account *openPayment.Account) error
}

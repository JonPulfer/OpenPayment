package account

import (
	"github.com/JonPulfer/OpenPayment/openPayment"
)

type Account interface {
	Create(openPayment.CreateAccountRequest) (openPayment.CreateAccountResponse, error)
}

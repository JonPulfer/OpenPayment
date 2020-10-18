package account

import (
	"context"

	"github.com/JonPulfer/OpenPayment/pkg/openPayment"
)

type Account interface {
	Create(ctx context.Context,
		request openPayment.CreateAccountRequest) (openPayment.CreateAccountResponse, error)
	Update(ctx context.Context,
		request openPayment.UpdateAccountRequest) (openPayment.UpdateAccountResponse, error)
	Balance(ctx context.Context,
		request openPayment.AccountBalanceRequest) (openPayment.AccountBalanceResponse, error)
}

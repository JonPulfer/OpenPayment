package account

import "github.com/JonPulfer/OpenPayment/pkg/openPayment/infrastructure/accountRepository"

type Account interface {
	Fetch(accountNumber int) (Account, error)
}

type SimpleAccount struct {
	accountRepository accountRepository.AccountRepository
}

func NewSimpleAccount(repos accountRepository.AccountRepository) *SimpleAccount {
	return &SimpleAccount{accountRepository: repos}
}

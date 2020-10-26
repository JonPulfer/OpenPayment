package account

type Account interface {
	Fetch(accountNumber int) (Account, error)
}

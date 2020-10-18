package openPayment

type Account struct {
	CustomerURLs []string `json:"customerURLs"`
	CardURLs     []string `json:"cardURLs"`
	Number       int      `json:"accountNumber"`
	SortCode     [3]int   `json:"sortCode"`
}

type CreateAccountRequest struct {
	CustomerURLs []string `json:"customerURLs"`
}

type CreateAccountResponse struct {
	AccountURL string `json:"accountURL"`
}

type UpdateAccountRequest struct {
	Changes map[string]interface{} `json:"changes"`
}

type UpdateAccountResponse struct {
	AccountURL string `json:"accountURL"`
}

type AccountBalanceRequest struct {
	AccountURL string `json:"accountURL"`
}

type AccountBalanceResponse struct {
	AccountURL string         `json:"accountURL"`
	Balance    AccountBalance `json:"accountBalance"`
}

type AccountBalance struct {
	Actual       float64 `json:"actual"`
	CurrencyCode string  `json:"currencyCode"`
}

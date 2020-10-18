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

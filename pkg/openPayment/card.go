package openPayment

import "time"

// Card registered to a Customer.
type Card struct {
	ID         string    `json:"-"`
	Number     [4][4]int `json:"cardNumber"`
	NameOnCard string    `json:"nameOnCard"`
	Expiry     time.Time `json:"expiry"`
	CVV        [3]int    `json:"-"`
}

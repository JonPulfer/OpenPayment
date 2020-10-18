package openPayment

import "time"

// Customer using this system.
type Customer struct {
	ShortName   string    `json:"shortName"`
	FullName    string    `json:"fullName"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Created     time.Time `json:"created"`
	Accounts    []string  `json:"accountURLs"`
	Cards       []string  `json:"cardURLs"`
}

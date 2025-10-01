package domain

import "fmt"

var validCurrencies = map[string]bool{
	"USD": true,
	"PEN": true,
}

type Money struct {
	Amount   int    `json:"amount"` // cents
	Currency string `json:"currency"`
}

func NewMoney(amount int, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, fmt.Errorf("amount cannot be negative")
	}

	if !validCurrencies[currency] {
		return Money{}, fmt.Errorf("invalid currency")
	}

	return Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

func (m Money) GetAmount() int {
	return m.Amount
}

func (m Money) GetCurrency() string {
	return m.Currency
}

package domain

import (
	"math/rand"
	"time"
)

type OrderCode string

func generateOrderCode() OrderCode {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return OrderCode(b)
}

package domain

import (
	"fmt"
	"regexp"
)

type Email string

func NewEmail(email string) (Email, error) {
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return "", fmt.Errorf("invalid email format")
	}

	return Email(email), nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Value() string {
	return string(e)
}

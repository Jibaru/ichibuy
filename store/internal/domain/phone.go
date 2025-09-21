package domain

import (
	"fmt"
	"regexp"
)

type Phone string

func NewPhone(phone string) (Phone, error) {
	if phone == "" {
		return "", fmt.Errorf("phone cannot be empty")
	}

	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		return "", fmt.Errorf("invalid phone format, should be +code number")
	}

	return Phone(phone), nil
}

func (p Phone) String() string {
	return string(p)
}

func (p Phone) Value() string {
	return string(p)
}

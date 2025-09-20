package domain

import "fmt"

type User struct {
	ID       string `sql:"id,primary"`
	Email    string `sql:"email"`
	Username string `sql:"username"`
}

func NewUser(id string, email string, username string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	return &User{
		ID:       id,
		Email:    email,
		Username: username,
	}, nil
}

func (u *User) TableName() string {
	return "users"
}

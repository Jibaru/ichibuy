package domain

import (
	"fmt"
	"time"
)

type Customer struct {
	ID        string    `sql:"id,primary"`
	FirstName string    `sql:"first_name"`
	LastName  string    `sql:"last_name"`
	Email     *string   `sql:"email"`
	Phone     *string   `sql:"phone"`
	UserID    string    `sql:"user_id"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}

func NewCustomer(id, firstName, lastName string, email *string, phone *string, userID string) (*Customer, error) {
	if firstName == "" {
		return nil, fmt.Errorf("firstName cannot be empty")
	}

	if len(firstName) > 50 {
		return nil, fmt.Errorf("firstName cannot exceed 50 characters")
	}

	if lastName == "" {
		return nil, fmt.Errorf("lastName cannot be empty")
	}

	if len(lastName) > 50 {
		return nil, fmt.Errorf("lastName cannot exceed 50 characters")
	}

	// Validate email if provided
	if email != nil {
		if _, err := NewEmail(*email); err != nil {
			return nil, err
		}
	}

	// Validate phone if provided
	if phone != nil {
		if _, err := NewPhone(*phone); err != nil {
			return nil, err
		}
	}

	now := time.Now().UTC()

	return &Customer{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Customer) Update(firstName, lastName string, email *string, phone *string, userID string) error {
	if userID != c.UserID {
		return fmt.Errorf("user id does not match")
	}

	if firstName == "" {
		return fmt.Errorf("firstName cannot be empty")
	}

	if len(firstName) > 50 {
		return fmt.Errorf("firstName cannot exceed 50 characters")
	}

	if lastName == "" {
		return fmt.Errorf("lastName cannot be empty")
	}

	if len(lastName) > 50 {
		return fmt.Errorf("lastName cannot exceed 50 characters")
	}

	// Validate email if provided
	if email != nil {
		if _, err := NewEmail(*email); err != nil {
			return err
		}
	}

	// Validate phone if provided
	if phone != nil {
		if _, err := NewPhone(*phone); err != nil {
			return err
		}
	}

	c.FirstName = firstName
	c.LastName = lastName
	c.Email = email
	c.Phone = phone
	c.UpdatedAt = time.Now().UTC()

	return nil
}

// Getters
func (c *Customer) GetID() string           { return c.ID }
func (c *Customer) GetFirstName() string    { return c.FirstName }
func (c *Customer) GetLastName() string     { return c.LastName }
func (c *Customer) GetUserID() string       { return c.UserID }
func (c *Customer) GetCreatedAt() time.Time { return c.CreatedAt }
func (c *Customer) GetUpdatedAt() time.Time { return c.UpdatedAt }

// GetEmail returns the email value object if set
func (c *Customer) GetEmail() *Email {
	if c.Email == nil {
		return nil
	}
	// We know this is valid since it was validated in NewCustomer/Update
	email, _ := NewEmail(*c.Email)
	return &email
}

// GetPhone returns the phone value object if set
func (c *Customer) GetPhone() *Phone {
	if c.Phone == nil {
		return nil
	}
	// We know this is valid since it was validated in NewCustomer/Update
	phone, _ := NewPhone(*c.Phone)
	return &phone
}

// GetEmailString returns the raw email string for database storage
func (c *Customer) GetEmailString() *string {
	return c.Email
}

// GetPhoneString returns the raw phone string for database storage
func (c *Customer) GetPhoneString() *string {
	return c.Phone
}

func (c *Customer) TableName() string {
	return "customers"
}

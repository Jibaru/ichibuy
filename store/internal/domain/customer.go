package domain

import (
	"encoding/json"
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

	Entity
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

	customer := &Customer{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	data, _ := json.Marshal(customer.createEventData())
	event := Event{
		ID:        fmt.Sprintf("%s_%v", customer.GetID(), customer.GetCreatedAt().Unix()),
		Type:      CustomerCreated,
		Data:      data,
		Timestamp: customer.GetCreatedAt(),
	}

	customer.events = append(customer.events, event)

	return customer, nil
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

	data, _ := json.Marshal(c.createEventData())
	event := Event{
		ID:        fmt.Sprintf("%s_%v", c.GetID(), c.GetUpdatedAt().Unix()),
		Type:      CustomerUpdated,
		Data:      data,
		Timestamp: c.GetUpdatedAt(),
	}

	c.events = append(c.events, event)

	return nil
}

func (c *Customer) createEventData() CustomerEventData {
	var emailStr, phoneStr *string
	if c.GetEmail() != nil {
		emailValue := c.GetEmail().Value()
		emailStr = &emailValue
	}
	if c.GetPhone() != nil {
		phoneValue := c.GetPhone().Value()
		phoneStr = &phoneValue
	}

	return CustomerEventData{
		ID:        c.GetID(),
		FirstName: c.GetFirstName(),
		LastName:  c.GetLastName(),
		Email:     emailStr,
		Phone:     phoneStr,
		UserID:    c.GetUserID(),
		CreatedAt: c.GetCreatedAt(),
		UpdatedAt: c.GetUpdatedAt(),
	}
}

func (c *Customer) PrepareDelete() {
	data, _ := json.Marshal(c.createEventData())

	event := Event{
		ID:        fmt.Sprintf("%s_%v_delete", c.GetID(), c.GetUpdatedAt().Unix()),
		Type:      CustomerDeleted,
		Data:      data,
		Timestamp: c.GetUpdatedAt(),
	}

	c.events = append(c.events, event)
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

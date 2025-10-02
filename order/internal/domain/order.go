package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Order struct {
	ID            string          `sql:"id,primary" json:"id"`
	Code          OrderCode       `sql:"code" json:"code"`
	CurrentStatus OrderStatus     `sql:"current_status" json:"current_status"`
	OrderLines    json.RawMessage `sql:"order_lines" json:"-order_lines"`
	CustomerID    string          `sql:"customer_id" json:"customer_id"`
	CreatedAt     time.Time       `sql:"created_at" json:"created_at"`
	UpdatedAt     time.Time       `sql:"updated_at" json:"updated_at"`

	orderLines []OrderLine
	Entity
}

func (o *Order) GetID() string                 { return o.ID }
func (o *Order) GetCode() OrderCode            { return o.Code }
func (o *Order) GetCurrentStatus() OrderStatus { return o.CurrentStatus }
func (o *Order) GetOrderLines() []OrderLine    { return o.orderLines }
func (o *Order) GetCustomerID() string         { return o.CustomerID }
func (o *Order) GetCreatedAt() time.Time       { return o.CreatedAt }

func (o *Order) Accept(acceptedBy string) error {
	switch o.CurrentStatus {
	case CreatedOrderStatus:
		o.CurrentStatus = AcceptedOrderStatus
		o.UpdatedAt = time.Now().UTC()
	default:
		return fmt.Errorf("order is not in created status")
	}
	return nil
}

func (o *Order) TableName() string {
	return "orders"
}

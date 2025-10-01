package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type OrderFactory struct {
	customerSvc CustomerService
	nextID      NextID
}

func NewOrderFactory(customerSvc CustomerService, nextID NextID) *OrderFactory {
	return &OrderFactory{
		customerSvc: customerSvc,
		nextID:      nextID,
	}

}

func (f *OrderFactory) NewOrder(
	ctx context.Context,
	orderLines []OrderLine,
	userID string,
) (*Order, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	if len(orderLines) == 0 {
		return nil, fmt.Errorf("orderLines cannot be empty")
	}

	// all order lines must have same currency
	currencies := map[string]bool{}
	storeIDs := map[string]bool{}
	for _, orderLine := range orderLines {
		currencies[orderLine.UnitPrice.GetCurrency()] = true
		storeIDs[orderLine.ProductStoreID] = true
	}

	if len(currencies) > 1 {
		return nil, fmt.Errorf("all order lines must have same currency")
	}

	if len(storeIDs) > 1 {
		return nil, fmt.Errorf("all order lines must have same store")
	}

	rawOrderLines, err := json.Marshal(orderLines)
	if err != nil {
		return nil, err
	}

	customer, err := f.customerSvc.FindByUserID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "find customer by user id failed", "error", err.Error())
		return nil, err
	}

	now := time.Now().UTC()

	order := &Order{
		ID:            f.nextID(),
		Code:          generateOrderCode(),
		CurrentStatus: CreatedOrderStatus,
		OrderLines:    rawOrderLines,
		CustomerID:    customer.ID,
		CreatedAt:     now,
		UpdatedAt:     now,

		orderLines: orderLines,
	}

	data, _ := json.Marshal(order)
	event := Event{
		ID:        fmt.Sprintf("%s_%v", order.GetID(), now.Unix()),
		Type:      OrderCreated,
		Data:      data,
		Timestamp: now,
	}

	order.events = append(order.events, event)

	return order, nil
}

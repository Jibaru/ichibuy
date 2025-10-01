package domain

import (
	"fmt"
	"strings"
)

type OrderLine struct {
	ID             string `json:"id"`
	ProductID      string `json:"product_id"`
	ProductName    string `json:"product_name"`
	ProductStoreID string `json:"product_store_id"`
	Quantity       int    `json:"quantity"`
	UnitPrice      Money  `json:"unit_price"`
}

func NewOrderLine(
	id string,
	productID string,
	productName string,
	productStoreID string,
	quantity int,
	unitPrice Money,
) (*OrderLine, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if strings.TrimSpace(productID) == "" {
		return nil, fmt.Errorf("productID cannot be empty")
	}

	if strings.TrimSpace(productName) == "" {
		return nil, fmt.Errorf("productName cannot be empty")
	}

	if strings.TrimSpace(productStoreID) == "" {
		return nil, fmt.Errorf("productStoreID cannot be empty")
	}

	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	if unitPrice.GetAmount() <= 0 {
		return nil, fmt.Errorf("unitPrice must be greater than 0")
	}

	return &OrderLine{
		ID:             id,
		ProductID:      productID,
		ProductName:    productName,
		ProductStoreID: productStoreID,
		Quantity:       quantity,
		UnitPrice:      unitPrice,
	}, nil
}

package domain

type OrderStatus string

const (
	CreatedOrderStatus  OrderStatus = "created"
	AcceptedOrderStatus OrderStatus = "accepted"
	FinishedOrderStatus OrderStatus = "finished"
	CanceledOrderStatus OrderStatus = "canceled"
	RejectedOrderStatus OrderStatus = "rejected"
)

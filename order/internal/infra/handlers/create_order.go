package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/order/internal/services"
)

type CreateOrderBody struct {
	OrderLines []OrderLineReq `json:"order_lines" binding:"required"`
}

type OrderLineReq struct {
	ProductID         string `json:"product_id" binding:"required"`
	ProductName       string `json:"product_name" binding:"required"`
	ProductStoreID    string `json:"product_store_id" binding:"required"`
	Quantity          int    `json:"quantity" binding:"required"`
	UnitPriceAmount   int    `json:"unit_price_amount" binding:"required"`
	UnitPriceCurrency string `json:"unit_price_currency" binding:"required"`
}

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Create a new order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order body CreateOrderBody true "Order data"
// @Success      201  {object}  services.CreateOrderResp
// @Failure      400  {object}  ErrorResp
// @Failure      401  {object}  ErrorResp
// @Router       /api/v1/orders [post]
// @Security     BearerAuth
func CreateOrder(CreateOrderService *services.CreateOrder) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		var req CreateOrderBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		resp, err := CreateOrderService.Exec(c, createOrderBodyToCreateOrderReq(req, userID.(string)))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}

func createOrderBodyToCreateOrderReq(body CreateOrderBody, userID string) services.CreateOrderReq {
	orderLines := []services.OrderLineReq{}
	for _, orderLine := range body.OrderLines {
		orderLines = append(orderLines, services.OrderLineReq{
			ProductID:         orderLine.ProductID,
			ProductName:       orderLine.ProductName,
			ProductStoreID:    orderLine.ProductStoreID,
			Quantity:          orderLine.Quantity,
			UnitPriceAmount:   orderLine.UnitPriceAmount,
			UnitPriceCurrency: orderLine.UnitPriceCurrency,
		})
	}
	return services.CreateOrderReq{
		OrderLines: orderLines,
		UserID:     userID,
	}
}

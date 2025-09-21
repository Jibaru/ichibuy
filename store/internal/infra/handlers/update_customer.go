package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

type UpdateCustomerBody struct {
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
}

// UpdateCustomer godoc
// @Summary      Update customer by ID
// @Description  Update a customer's information
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Param        customer body UpdateCustomerBody true "Customer data"
// @Success      204
// @Failure      400  {object}  ErrorResp
// @Router       /api/v1/customers/{id} [put]
// @Security     BearerAuth
func UpdateCustomer(updateCustomerService *services.UpdateCustomer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "id parameter is required"})
			return
		}

		var req UpdateCustomerBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		err := updateCustomerService.Exec(c, services.UpdateCustomerReq{
			ID:        id,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
			UserID:    userID.(string),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

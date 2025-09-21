package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

type CreateCustomerBody struct {
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
}

// CreateCustomer godoc
// @Summary      Create a new customer
// @Description  Create a new customer with first name, last name, email and phone
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customer body CreateCustomerBody true "Customer data"
// @Success      201  {object}  services.CreateCustomerResp
// @Failure      400  {object}  ErrorResp
// @Failure      401  {object}  ErrorResp
// @Router       /api/v1/customers [post]
// @Security     BearerAuth
func CreateCustomer(createCustomerService *services.CreateCustomer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		var req CreateCustomerBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		resp, err := createCustomerService.Exec(c, services.CreateCustomerReq{
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

		c.JSON(http.StatusCreated, resp)
	}
}

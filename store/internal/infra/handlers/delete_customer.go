package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// DeleteCustomer godoc
// @Summary      Delete customer by ID
// @Description  Delete a customer
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Success      204
// @Failure      400  {object}  ErrorResp
// @Router       /api/v1/customers/{id} [delete]
// @Security     BearerAuth
func DeleteCustomer(deleteCustomerService *services.DeleteCustomer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "id parameter is required"})
			return
		}

		err := deleteCustomerService.Exec(c, services.DeleteCustomerReq{ID: id})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

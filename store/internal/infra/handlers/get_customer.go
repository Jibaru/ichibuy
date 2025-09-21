package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// GetCustomer godoc
// @Summary      Get customer by ID
// @Description  Retrieve a customer by its ID
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Success      200  {object}  services.GetCustomerResp
// @Failure      400  {object}  ErrorResp
// @Failure      404  {object}  ErrorResp
// @Router       /api/v1/customers/{id} [get]
// @Security     BearerAuth
func GetCustomer(getCustomerService *services.GetCustomer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "id parameter is required"})
			return
		}

		resp, err := getCustomerService.Exec(c, services.GetCustomerReq{ID: id})
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResp{Error: "customer not found"})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

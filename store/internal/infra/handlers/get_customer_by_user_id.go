package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// GetCustomerByUserID godoc
// @Summary      Get customer by user ID
// @Description  Retrieve a customer using the associated user ID
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        userId path string true "User ID"
// @Success      200  {object}  services.GetCustomerByUserIDResp
// @Failure      400  {object}  ErrorResp
// @Failure      404  {object}  ErrorResp
// @Router       /api/v1/customers/user/{userId} [get]
// @Security     BearerAuth
func GetCustomerByUserID(service *services.GetCustomerByUserID) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("userId")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "userId parameter is required"})
			return
		}

		resp, err := service.Exec(c, services.GetCustomerByUserIDReq{UserID: id})
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResp{Error: "customer not found"})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

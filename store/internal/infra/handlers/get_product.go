package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// GetProduct godoc
// @Summary      Get a product by ID
// @Description  Get a single product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  services.GetProductResp
// @Failure      400  {object}  ErrorResp
// @Failure      401  {object}  ErrorResp
// @Failure      404  {object}  ErrorResp
// @Router       /api/v1/products/{id} [get]
// @Security     BearerAuth
func GetProduct(getProductService *services.GetProduct) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "id is required"})
			return
		}

		resp, err := getProductService.Exec(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
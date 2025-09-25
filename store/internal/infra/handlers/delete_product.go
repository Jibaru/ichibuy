package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// DeleteProduct godoc
// @Summary      Delete product by ID
// @Description  Delete a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      204
// @Failure      400  {object}  ErrorResp
// @Failure      401  {object}  ErrorResp
// @Failure      404  {object}  ErrorResp
// @Router       /api/v1/products/{id} [delete]
// @Security     BearerAuth
func DeleteProduct(deleteProductService *services.DeleteProduct) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "id parameter is required"})
			return
		}

		err := deleteProductService.Exec(c, services.DeleteProductReq{ID: id})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
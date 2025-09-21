package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// DeleteStore godoc
// @Summary      Delete store by ID
// @Description  Delete a store
// @Tags         stores
// @Accept       json
// @Produce      json
// @Param        id path string true "Store ID"
// @Success      204
// @Failure      400  {object}  ErrorResp
// @Router       /api/v1/stores/{id} [delete]
// @Security     BearerAuth
func DeleteStore(deleteStoreService *services.DeleteStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "id parameter is required"})
			return
		}

		err := deleteStoreService.Exec(c, services.DeleteStoreReq{ID: id})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// GetStore godoc
// @Summary      Get store by ID
// @Description  Retrieve a store by its ID
// @Tags         stores
// @Accept       json
// @Produce      json
// @Param        id path string true "Store ID"
// @Success      200  {object}  services.GetStoreResp
// @Failure      400  {object}  ErrorResp
// @Failure      404  {object}  ErrorResp
// @Router       /api/v1/stores/{id} [get]
// @Security     BearerAuth
func GetStore(getStoreService *services.GetStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "id parameter is required"})
			return
		}

		store, err := getStoreService.Exec(c, services.GetStoreReq{ID: id})
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResp{Error: "store not found"})
			return
		}

		c.JSON(http.StatusOK, store)
	}
}

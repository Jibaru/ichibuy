package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/services"
)

type UpdateStoreBody struct {
	Name        string          `json:"name" binding:"required"`
	Description *string         `json:"description"`
	Location    domain.Location `json:"location" binding:"required"`
}

// UpdateStore godoc
// @Summary      Update store by ID
// @Description  Update a store's information
// @Tags         stores
// @Accept       json
// @Produce      json
// @Param        id path string true "Store ID"
// @Param        store body UpdateStoreBody true "Store data"
// @Success      204
// @Failure      400  {object}  ErrorResp
// @Router       /api/v1/stores/{id} [put]
// @Security     BearerAuth
func UpdateStore(updateStoreService *services.UpdateStore) gin.HandlerFunc {
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

		var req UpdateStoreBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		err := updateStoreService.Exec(c, services.UpdateStoreReq{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
			Location:    req.Location,
			UserID:      userID.(string),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

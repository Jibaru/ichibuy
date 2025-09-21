package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/services"
)

type CreateStoreBody struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Lat         float64 `json:"lat" binding:"required"`
	Lng         float64 `json:"lng" binding:"required"`
}

// CreateStore godoc
// @Summary      Create a new store
// @Description  Create a new store with name, description and location
// @Tags         stores
// @Accept       json
// @Produce      json
// @Param        store body CreateStoreBody true "Store data"
// @Success      201  {object}  services.CreateStoreResp
// @Failure      400  {object}  ErrorResp
// @Failure      401  {object}  ErrorResp
// @Router       /api/v1/stores [post]
// @Security     BearerAuth
func CreateStore(createStoreService *services.CreateStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		var req CreateStoreBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		resp, err := createStoreService.Exec(c, services.CreateStoreReq{
			Name:        req.Name,
			Description: req.Description,
			Location:    domain.Location{Lat: req.Lat, Lng: req.Lng},
			UserID:      userID.(string),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}

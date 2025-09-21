package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// ListStores godoc
// @Summary      List stores
// @Description  Get paginated list of stores with filters and sorting
// @Tags         stores
// @Accept       json
// @Produce      json
// @Param        name query string false "Filter by name"
// @Param        description query string false "Filter by description"
// @Param        sort_by query string false "Sort by field" default("name")
// @Param        sort_order query string false "Sort order" default("ASC")
// @Param        offset query int false "Offset" default(0)
// @Param        limit query int false "Limit" default(10)
// @Success      200  {object}  services.ListStoresResp
// @Failure      401  {object}  ErrorResp
// @Failure      500  {object}  ErrorResp
// @Router       /api/v1/stores [get]
// @Security     BearerAuth
func ListStores(listStoresService *services.ListStores) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		filters := services.StoreFilters{
			UserID: userID.(string),
		}

		if name := c.Query("name"); name != "" {
			filters.Name = &name
		}

		if description := c.Query("description"); description != "" {
			filters.Description = &description
		}

		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		pagination := services.Pagination{
			Offset: offset,
			Limit:  limit,
		}

		sorting := services.Sorting{
			Field: c.DefaultQuery("sort_by", "name"),
			Order: c.DefaultQuery("sort_order", "ASC"),
		}

		serviceReq := services.ListStoresReq{
			Filters:    filters,
			Pagination: pagination,
			Sorting:    sorting,
		}

		resp, err := listStoresService.Exec(c, serviceReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"stores": resp.Stores,
			"total":  resp.Total,
			"offset": offset,
			"limit":  limit,
		})
	}
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// ListProducts godoc
// @Summary      List products
// @Description  Get paginated list of products with filters and sorting
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        store_id query string false "Filter by store ID"
// @Param        name query string false "Filter by name"
// @Param        description query string false "Filter by description"
// @Param        active query bool false "Filter by active status"
// @Param        sort_by query string false "Sort by field" default("name")
// @Param        sort_order query string false "Sort order" default("ASC")
// @Param        offset query int false "Offset" default(0)
// @Param        limit query int false "Limit" default(10)
// @Success      200  {object}  services.ListProductsResp
// @Failure      401  {object}  ErrorResp
// @Failure      500  {object}  ErrorResp
// @Router       /api/v1/products [get]
// @Security     BearerAuth
func ListProducts(listProductsService *services.ListProducts) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		filters := services.ProductFilters{}

		if storeID := c.Query("store_id"); storeID != "" {
			filters.StoreID = storeID
		}

		if name := c.Query("name"); name != "" {
			filters.Name = &name
		}

		if description := c.Query("description"); description != "" {
			filters.Description = &description
		}

		if activeStr := c.Query("active"); activeStr != "" {
			if active, err := strconv.ParseBool(activeStr); err == nil {
				filters.Active = &active
			}
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

		serviceReq := services.ListProductsReq{
			Filters:    filters,
			Pagination: pagination,
			Sorting:    sorting,
		}

		resp, err := listProductsService.Exec(c, serviceReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"products": resp.Products,
			"total":    resp.Total,
			"offset":   offset,
			"limit":    limit,
		})
	}
}
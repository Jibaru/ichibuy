package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
	"ichibuy/store/internal/shared/context"
)

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create a new product with name, description, active status, store ID, images (as files) and prices
// @Tags         products
// @Accept       multipart/form-data
// @Produce      json
// @Param        name formData string true "Product name"
// @Param        description formData string false "Product description"
// @Param        active formData bool true "Product active status"
// @Param        store_id formData string true "Store ID"
// @Param        prices formData string true "JSON array of prices"
// @Param        images formData file false "Product images (multiple files allowed)"
// @Success      201  {object}  services.CreateProductResp
// @Failure      400  {object}  ErrorResp
// @Failure      401  {object}  ErrorResp
// @Router       /api/v1/products [post]
// @Security     BearerAuth
func CreateProduct(createProductService *services.CreateProduct) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.InfoContext(c, "token", "value", c.Value(context.APITokenKey))

		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "failed to parse multipart form: " + err.Error()})
			return
		}

		// Get form values
		name := c.PostForm("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "name is required"})
			return
		}

		storeID := c.PostForm("store_id")
		if storeID == "" {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "store_id is required"})
			return
		}

		active, _ := strconv.ParseBool(c.DefaultPostForm("active", "true"))

		var description *string
		if desc := c.PostForm("description"); desc != "" {
			description = &desc
		}

		// Parse prices JSON
		var prices []NewPriceDTO
		if pricesStr := c.PostForm("prices"); pricesStr != "" {
			if err := json.Unmarshal([]byte(pricesStr), &prices); err != nil {
				c.JSON(http.StatusBadRequest, ErrorResp{Error: "invalid prices JSON: " + err.Error()})
				return
			}
		}

		// Get uploaded image files
		form, _ := c.MultipartForm()
		imageFiles := form.File["images"]

		// Convert multipart files to FileDTO
		fileDTOs, err := convertMultipartFilesToDTOs(imageFiles)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: "failed to process uploaded files: " + err.Error()})
			return
		}

		resp, err := createProductService.Exec(c, services.CreateProductReq{
			Name:        name,
			Description: description,
			Active:      active,
			StoreID:     storeID,
			ImageFiles:  fileDTOs,
			Prices:      convertHandlerPriceDTOsToService(prices),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resp)
	}
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ichibuy/store/internal/services"
)

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update an existing product. Note: StoreID, ID and CreatedAt cannot be updated. Images and Prices are replaced entirely.
// @Tags         products
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path string true "Product ID"
// @Param        name formData string true "Product name"
// @Param        description formData string false "Product description"
// @Param        active formData bool true "Product active status"
// @Param        prices formData string true "JSON array of prices"
// @Param        deleteImagesIDs formData string false "JSON array of image IDs to delete"
// @Param        deletePricesIDs formData string false "JSON array of price IDs to delete"
// @Param        images formData file false "Product images (multiple files allowed)"
// @Success      204
// @Failure      400     {object} ErrorResp
// @Failure      401     {object} ErrorResp
// @Failure      404     {object} ErrorResp
// @Router       /api/v1/products/{id} [put]
// @Security     BearerAuth
func UpdateProduct(updateProductService *services.UpdateProduct) gin.HandlerFunc {
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

		// Parse delete IDs
		var deleteImageIDs []string
		if deleteImagesStr := c.PostForm("deleteImagesIDs"); deleteImagesStr != "" {
			if err := json.Unmarshal([]byte(deleteImagesStr), &deleteImageIDs); err != nil {
				c.JSON(http.StatusBadRequest, ErrorResp{Error: "invalid deleteImagesIDs JSON: " + err.Error()})
				return
			}
		}

		var deletePricesIDs []string
		if deletePricesStr := c.PostForm("deletePricesIDs"); deletePricesStr != "" {
			if err := json.Unmarshal([]byte(deletePricesStr), &deletePricesIDs); err != nil {
				c.JSON(http.StatusBadRequest, ErrorResp{Error: "invalid deletePricesIDs JSON: " + err.Error()})
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

		err = updateProductService.Exec(c, services.UpdateProductReq{
			ID:              id,
			Name:            name,
			Description:     description,
			Active:          active,
			NewImageFiles:   fileDTOs,
			NewPrices:       convertHandlerPriceDTOsToService(prices),
			DeleteImageIDs:  deleteImageIDs,
			DeletePricesIDs: deletePricesIDs,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

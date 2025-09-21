package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/auth/internal/services"
)

// StartOAuth godoc
// @Summary      StartOAuth
// @Description  StartOAuth
// @Accept       json
// @Produce      json
// @Param        provider path string true "Provider: google"
// @Success      307
// @Failure      500    {object}    ErrorResp
// @Router       /api/v1/auth/{provider} [get]
func StartOAuth(startOAuth *services.StartOAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := startOAuth.Exec(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResp{Error: err.Error()})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, resp.URL)
	}
}

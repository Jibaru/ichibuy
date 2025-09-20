package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ichibuy/auth/internal/services"
)

// OAuthCallback godoc
// @Summary      OAuthCallback
// @Description  OAuthCallback
// @Accept       json
// @Produce      json
// @Success      307
// @Failure      400    {object}    ErrorResp
// @Router       /api/v1/auth/{provider}/callback [get]
func OAuthCallback(finishOAuth *services.FinishOAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := finishOAuth.Exec(c, services.FinishOAuthReq{
			Code: c.Query("code"),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, resp.URL)
	}
}

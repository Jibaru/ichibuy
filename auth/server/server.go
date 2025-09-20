package server

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"ichibuy/auth/config"
	"ichibuy/auth/internal/infra/handlers"
	"ichibuy/auth/internal/infra/middlewares"
	"ichibuy/auth/internal/infra/persistence/postgres"
	infraServices "ichibuy/auth/internal/infra/services"
	"ichibuy/auth/internal/services"
)

func New(cfg config.Config, db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.UseCORS())

	googleOAuthConfig := &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/api/v1/auth/google/callback", cfg.APIBaseURI),
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	userDAO := postgres.NewUserDAO(db)

	nextIDFunc := uuid.NewString

	startOAuthServ := services.NewStartOAuth(googleOAuthConfig)
	finishOAuthServ := services.NewFinishOAuth(userDAO, googleOAuthConfig, infraServices.GoogleInfoExtractor, nextIDFunc, cfg)

	api := router.Group("/api/v1")
	{
		api.GET("/auth/google", handlers.StartOAuth(startOAuthServ))
		api.GET("/auth/google/callback", handlers.OAuthCallback(finishOAuthServ))
		api.GET("/auth/.well-known/jwks.json", handlers.GetJWKS(cfg))
	}

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

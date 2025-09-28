package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	authHTTP "github.com/Jibaru/ichibuy/api-client/go/auth"

	"ichibuy/order/config"
	"ichibuy/order/internal/infra/events"
	"ichibuy/order/internal/infra/middlewares"
	"ichibuy/order/internal/infra/persistence/postgres"
)

func New(cfg config.Config, db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.UseCORS())

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	authClient := authHTTP.NewAPIClient(&authHTTP.Configuration{
		BasePath:   cfg.AuthBaseURL,
		HTTPClient: httpClient,
	})

	jwtMiddleware := middlewares.NewJWTAuthMiddleware(authClient)

	// DAOs
	eventDAO := postgres.NewEventDAO(db)

	eventBus := events.NewBus(eventDAO)
	nextIDFunc := uuid.NewString

	fmt.Println(eventBus, nextIDFunc)

	// Use-Cases

	// Routes
	api := router.Group("/api/v1")
	api.Use(jwtMiddleware.ValidateToken())
	{
		orders := api.Group("/products")
		{
			fmt.Println(orders)
		}
	}

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

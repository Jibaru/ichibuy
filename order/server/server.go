package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	authHTTP "github.com/Jibaru/ichibuy/api-client/go/auth"
	storeHTTP "github.com/Jibaru/ichibuy/api-client/go/store"

	"ichibuy/order/config"
	"ichibuy/order/internal/domain"
	"ichibuy/order/internal/infra/events"
	"ichibuy/order/internal/infra/handlers"
	"ichibuy/order/internal/infra/middlewares"
	"ichibuy/order/internal/infra/persistence/postgres"
	infraServices "ichibuy/order/internal/infra/services"
	"ichibuy/order/internal/services"
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
	storeClient := storeHTTP.NewAPIClient(&storeHTTP.Configuration{
		BasePath:   cfg.StoreBaseURL,
		HTTPClient: httpClient,
	})

	jwtMiddleware := middlewares.NewJWTAuthMiddleware(authClient)

	// DAOs
	eventDAO := postgres.NewEventDAO(db)
	orderDAO := postgres.NewOrderDAO(db)

	eventBus := events.NewBus(eventDAO)
	nextIDFunc := uuid.NewString

	// Domain Services
	customerSvc := infraServices.NewCustomerService(storeClient)

	// Factories
	orderFactory := domain.NewOrderFactory(customerSvc, nextIDFunc)

	// Use-Cases
	createOrderService := services.NewCreateOrder(orderDAO, eventBus, nextIDFunc, orderFactory)

	// Routes
	api := router.Group("/api/v1")
	api.Use(jwtMiddleware.ValidateToken())
	{
		orders := api.Group("/orders")
		{
			orders.POST("", handlers.CreateOrder(createOrderService))
			// list my orders
			// cancel order (by me)
			// accept order (by the store owner)
			// reject order (by the store owner)

		}
	}

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

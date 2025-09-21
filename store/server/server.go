package server

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ichibuy/store/config"
	"ichibuy/store/internal/infra/client"
	"ichibuy/store/internal/infra/events"
	"ichibuy/store/internal/infra/handlers"
	"ichibuy/store/internal/infra/middlewares"
	"ichibuy/store/internal/infra/persistence/postgres"
	"ichibuy/store/internal/services"
)

func New(cfg config.Config, db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.UseCORS())

	authClient := client.NewAuthClient(cfg.AuthBaseURL)
	jwtMiddleware := middlewares.NewJWTAuthMiddleware(authClient)

	storeDAO := postgres.NewStoreDAO(db)
	customerDAO := postgres.NewCustomerDAO(db)

	eventBus := events.NewMemoryEventBus()
	nextIDFunc := uuid.NewString

	createStoreService := services.NewCreateStore(storeDAO, eventBus, nextIDFunc)
	getStoreService := services.NewGetStore(storeDAO)
	updateStoreService := services.NewUpdateStore(storeDAO, eventBus, nextIDFunc)
	deleteStoreService := services.NewDeleteStore(storeDAO, eventBus, nextIDFunc)
	listStoresService := services.NewListStores(storeDAO)

	createCustomerService := services.NewCreateCustomer(customerDAO, eventBus, nextIDFunc)
	getCustomerService := services.NewGetCustomer(customerDAO)
	updateCustomerService := services.NewUpdateCustomer(customerDAO, eventBus, nextIDFunc)
	deleteCustomerService := services.NewDeleteCustomer(customerDAO, eventBus, nextIDFunc)

	api := router.Group("/api/v1")
	api.Use(jwtMiddleware.ValidateToken())
	{
		stores := api.Group("/stores")
		{
			stores.POST("", handlers.CreateStore(createStoreService))
			stores.GET("/:id", handlers.GetStore(getStoreService))
			stores.PUT("/:id", handlers.UpdateStore(updateStoreService))
			stores.DELETE("/:id", handlers.DeleteStore(deleteStoreService))
			stores.GET("", handlers.ListStores(listStoresService))
		}

		customers := api.Group("/customers")
		{
			customers.POST("", handlers.CreateCustomer(createCustomerService))
			customers.GET("/:id", handlers.GetCustomer(getCustomerService))
			customers.PUT("/:id", handlers.UpdateCustomer(updateCustomerService))
			customers.DELETE("/:id", handlers.DeleteCustomer(deleteCustomerService))
		}

		api.POST("/graphql", handlers.GraphQLStores(listStoresService))
	}

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

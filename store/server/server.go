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
	fstorageHTTP "github.com/Jibaru/ichibuy/api-client/go/fstorage"

	"ichibuy/store/config"
	"ichibuy/store/internal/infra/events"
	"ichibuy/store/internal/infra/handlers"
	"ichibuy/store/internal/infra/middlewares"
	"ichibuy/store/internal/infra/persistence/postgres"
	infraServices "ichibuy/store/internal/infra/services"
	"ichibuy/store/internal/services"
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
	fstorageClient := fstorageHTTP.NewAPIClient(&fstorageHTTP.Configuration{
		BasePath:   cfg.FStorageBaseURL,
		HTTPClient: httpClient,
	})

	jwtMiddleware := middlewares.NewJWTAuthMiddleware(authClient)

	storeDAO := postgres.NewStoreDAO(db)
	customerDAO := postgres.NewCustomerDAO(db)
	productDAO := postgres.NewProductDAO(db)

	eventBus := events.NewMemoryEventBus()
	nextIDFunc := uuid.NewString

	storageSvc := infraServices.NewStorageService(fstorageClient)

	createStoreService := services.NewCreateStore(storeDAO, eventBus, nextIDFunc)
	getStoreService := services.NewGetStore(storeDAO)
	updateStoreService := services.NewUpdateStore(storeDAO, eventBus, nextIDFunc)
	deleteStoreService := services.NewDeleteStore(storeDAO, eventBus, nextIDFunc)
	listStoresService := services.NewListStores(storeDAO)

	createCustomerService := services.NewCreateCustomer(customerDAO, eventBus, nextIDFunc)
	getCustomerService := services.NewGetCustomer(customerDAO)
	updateCustomerService := services.NewUpdateCustomer(customerDAO, eventBus, nextIDFunc)
	deleteCustomerService := services.NewDeleteCustomer(customerDAO, eventBus, nextIDFunc)

	createProductService := services.NewCreateProduct(productDAO, eventBus, nextIDFunc, storageSvc)
	getProductService := services.NewGetProduct(productDAO)
	updateProductService := services.NewUpdateProduct(productDAO, eventBus, nextIDFunc, storageSvc)
	deleteProductService := services.NewDeleteProduct(productDAO, eventBus, nextIDFunc)
	listProductsService := services.NewListProducts(productDAO)

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

		products := api.Group("/products")
		{
			products.POST("", handlers.CreateProduct(createProductService))
			products.GET("/:id", handlers.GetProduct(getProductService))
			products.PUT("/:id", handlers.UpdateProduct(updateProductService))
			products.DELETE("/:id", handlers.DeleteProduct(deleteProductService))
			products.GET("", handlers.ListProducts(listProductsService))
		}

		api.POST("/graphql", handlers.GraphQLStores(listStoresService))
	}

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

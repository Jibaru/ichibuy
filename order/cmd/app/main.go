package main

import (
	"ichibuy/order/config"
	"ichibuy/order/db"
	_ "ichibuy/order/docs"
	"ichibuy/order/server"
)

// @title           ichibuy/order API
// @version         1.0
// @description     This is the ichibuy/order API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      ichibuy-order.vercel.app

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cfg := config.Load()
	db, err := db.New(cfg.PostgresURI)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := server.New(cfg, db)
	router.Run(":" + cfg.APIPort)
}

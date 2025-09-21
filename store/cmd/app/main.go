package main

import (
	"ichibuy/store/config"
	"ichibuy/store/db"
	_ "ichibuy/store/docs"
	"ichibuy/store/server"
)

// @title           ichibuy/store API
// @version         1.0
// @description     This is the ichibuy/store API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      https://ichibuy-store.vercel.app

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
